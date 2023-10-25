package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/barpav/demography/internal/rest/models"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/demography-api/#/people/get_people
func (s *Service) searchByData(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", models.MimeTypeSearchResultV1:
		s.searchByDataV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) searchByDataV1(w http.ResponseWriter, r *http.Request) {
	filters, err := searchFilters(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var result *models.SearchResultV1
	result, err = s.storage.SearchResultV1(r.Context(), filters)

	if err != nil {
		log.Err(err).Msg("Failed to receive search result (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", models.MimeTypeSearchResultV1)
	err = json.NewEncoder(w).Encode(result)

	if err != nil {
		log.Err(err).Msg("Failed to serialize search result (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info().Msg(fmt.Sprintf("Search results: %d", result.Total))
}

func searchFilters(r *http.Request) (filters *models.SearchFilters, err error) {
	query := r.URL.Query()
	filters = &models.SearchFilters{
		Surname:    query.Get("surname"),
		Name:       query.Get("name"),
		Patronymic: query.Get("patronymic"),
		Gender:     query.Get("gender"),
		Country:    query.Get("country"),
	}

	var parseErr error
	var param int64
	param, parseErr = integerQueryParameter(r, "age")

	if parseErr != nil {
		err = errors.Join(err, parseErr)
	} else {
		filters.Age = int(param)
	}

	param, parseErr = integerQueryParameter(r, "after")

	if parseErr != nil {
		err = errors.Join(err, parseErr)
	} else {
		filters.After = param
	}

	const limitDefault = 30

	if query.Get("limit") == "" {
		filters.Limit = limitDefault
	} else {
		const limitMin = 1
		const limitMax = 100

		param, parseErr = integerQueryParameter(r, "limit")

		if parseErr != nil {
			err = errors.Join(err, parseErr)
		} else {
			if param < limitMin || param > limitMax {
				err = errors.Join(err, fmt.Errorf("Invalid parameter 'limit': min %d, max %d.", limitMin, limitMax))
			} else {
				filters.Limit = int(param)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return filters, nil
}

func integerQueryParameter(r *http.Request, name string) (value int64, err error) {
	param := r.URL.Query().Get(name)

	if param == "" {
		return 0, nil
	}

	value, err = strconv.ParseInt(param, 10, 0)

	if err != nil {
		return 0, fmt.Errorf("Parameter '%s' must be an integer type.", name)
	}

	return value, nil
}
