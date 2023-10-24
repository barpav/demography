package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/barpav/demography/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/demography-api/#/people/get_people__id_
func (s *Service) getPersonData(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", models.MimeTypeEnrichedPersonDataV1:
		s.getPersonDataV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) getPersonDataV1(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var data *models.EnrichedPersonDataV1
	data, err = s.storage.EnrichedPersonDataV1(r.Context(), id)

	if err != nil {
		log.Err(err).Msg("Failed to receive enriched person data (v1) by id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", models.MimeTypeEnrichedPersonDataV1)
	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Err(err).Msg("Failed to serialize enriched person data (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
