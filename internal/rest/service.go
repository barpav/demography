package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	shutdown chan struct{}
	cfg      *config
	server   *http.Server
	storage  Storage
}

type Storage interface{}

func (s *Service) Start(storage Storage) {
	s.cfg = &config{}
	s.cfg.Read()

	s.storage = storage

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.port),
		Handler: s.operations(),
	}

	s.shutdown = make(chan struct{}, 1)

	go func() {
		err := s.server.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Err(err).Msg("HTTP server crashed.")
		}

		s.shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	err = s.server.Shutdown(ctx)

	if err != nil {
		err = fmt.Errorf("failed to stop HTTP service: %w", err)
	}

	return err
}

func (s *Service) Shutdown() <-chan struct{} {
	return s.shutdown
}

// Specification: https://barpav.github.io/demography-api/#/people
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	// ops.Use(s.enableCORS)

	// ops.Post("/v1/people", s.addNewPerson)
	// ops.Get("/v1/people", s.searchByData)
	// ops.Get("/v1/people/{id}", s.getPersonData)
	// ops.Put("/v1/people/{id}", s.editPersonData)
	// ops.Delete("/v1/people/{id}", s.deletePersonData)

	return ops
}
