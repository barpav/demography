package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/demography-api/#/people/delete_people__id_
func (s *Service) deletePersonData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = s.storage.DeletePersonData(r.Context(), id)

	if err != nil {
		if _, ok := err.(ErrPersonDataNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Err(err).Msg("Failed to delete person data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
