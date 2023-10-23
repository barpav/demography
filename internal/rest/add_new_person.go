package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/barpav/demography/internal/rest/models"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/demography-api/#/people/post_people
func (s *Service) addNewPerson(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case models.MimeTypeNewPersonDataV1:
		s.addNewPersonV1(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func (s *Service) addNewPersonV1(w http.ResponseWriter, r *http.Request) {
	personData := models.NewPersonDataV1{}
	err := personData.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var fullData *models.EnrichedPersonDataV1
	fullData, err = s.enrichedPersonDataV1(&personData)

	if err != nil {
		log.Err(err).Msg("Failed to receive enriched person data (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.storage.CreateNewPersonDataV1(r.Context(), fullData)

	if err != nil {
		log.Err(err).Msg("Failed to save new person data (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", models.MimeTypeEnrichedPersonDataV1)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(fullData)

	if err != nil {
		log.Err(err).Msg("Failed to serialize enriched person data (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) enrichedPersonDataV1(data *models.NewPersonDataV1) (result *models.EnrichedPersonDataV1, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(s.cfg.statsTimeout))
	defer cancel()

	var age int
	var gender, country string

	const servicesTotal = 3

	wg := &sync.WaitGroup{}
	wg.Add(servicesTotal)
	done := make(chan struct{}, 1)
	interrupt := make(chan struct{}, servicesTotal)

	// waiting for confirmation from all 3rd parties (all or nothing)
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	// receiving age statistics from 3rd party (retries in timeout range)
	go func() {
		var statsErr error
		for {
			select {
			case <-interrupt:
				wg.Done()
				return
			default:
				age, statsErr = s.stats.AgeByName(data.Name)
				if statsErr == nil {
					wg.Done()
					return
				} else {
					log.Err(statsErr).Msg("Failed to receive age statistics.")
				}
			}
		}
	}()

	// receiving gender statistics from 3rd party (retries in timeout range)
	go func() {
		var statsErr error
		for {
			select {
			case <-interrupt:
				wg.Done()
				return
			default:
				gender, statsErr = s.stats.GenderByName(data.Name)
				if statsErr == nil {
					wg.Done()
					return
				} else {
					log.Err(statsErr).Msg("Failed to receive gender statistics.")
				}
			}
		}
	}()

	// receiving country statistics from 3rd party (retries in timeout range)
	go func() {
		var statsErr error
		for {
			select {
			case <-interrupt:
				wg.Done()
				return
			default:
				country, statsErr = s.stats.CountryByName(data.Name)
				if statsErr == nil {
					wg.Done()
					return
				} else {
					log.Err(statsErr).Msg("Failed to receive country statistics.")
				}
			}
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		for i := servicesTotal; i > 0; i-- {
			interrupt <- struct{}{}
		}
		return nil, fmt.Errorf("failed to enrich person data: %w", ctx.Err())
	}

	result = &models.EnrichedPersonDataV1{
		Surname:    data.Surname,
		Name:       data.Name,
		Patronymic: data.Patronymic,
		Age:        age,
		Gender:     gender,
		Country:    country,
	}

	return result, nil
}
