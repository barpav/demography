package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/barpav/demography/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/demography-api/#/people/put_people__id_
func (s *Service) editPersonData(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case models.MimeTypeEditedPersonDataV1:
		s.editPersonDataV1(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func (s *Service) editPersonDataV1(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	editedData := models.EditedPersonDataV1{}
	err = editedData.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = s.storage.UpdatePersonDataV1(r.Context(), id, &editedData)

	if err != nil {
		if _, ok := err.(ErrPersonDataNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Err(err).Msg("Failed to update person data (v1).")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info().Msg(fmt.Sprintf("Person data with id '%d' edited.", id))
}

type ErrPersonDataNotFound interface {
	Error() string
	ImplementsPersonDataNotFoundError()
}
