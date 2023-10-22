package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/barpav/demography/internal/data"
	"github.com/barpav/demography/internal/rest"
)

func main() {
	app := microservice{}
	err := app.launch()

	if err == nil {
		log.Info().Msg("Microservice launched.")
	} else {
		log.Err(err).Msg("Failed to launch microservice.")
		app.abort()
	}

	err = app.serveAndShutdownGracefully()

	if err == nil {
		log.Info().Msg("Microservice stopped.")
	} else {
		log.Err(err).Msg("Failed to shutdown microservice gracefully.")
	}
}

type microservice struct {
	api struct {
		public *rest.Service // specification: https://barpav.github.io/demography-api/#/people
	}
	storage  *data.Storage
	shutdown chan os.Signal
}

func (m *microservice) launch() (err error) {
	m.shutdown = make(chan os.Signal, 2)
	signal.Notify(m.shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	m.storage = &data.Storage{}
	err = m.storage.Open()

	m.api.public = &rest.Service{}
	m.api.public.Start(m.storage)

	return err
}

func (m *microservice) abort() {
	m.shutdown <- syscall.SIGINT
}

func (m *microservice) serveAndShutdownGracefully() (err error) {
	select {
	case <-m.shutdown:
	case <-m.api.public.Shutdown():
	}

	log.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = errors.Join(err, m.api.public.Stop(ctx))
	err = errors.Join(err, m.storage.Close(ctx))

	return err
}
