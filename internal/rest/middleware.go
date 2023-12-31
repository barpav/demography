package rest

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Normally done by front end web server.
func (s *Service) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg(fmt.Sprintf("%s %s", r.Method, r.RequestURI))

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, HEAD, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", "1728000") // valid for 20 days
			w.WriteHeader(http.StatusNoContent)
			log.Debug().Msg("CORS options returned")
			return
		}

		w.Header().Set("Access-Control-Expose-Headers",
			"Content-Type,Content-Length,Content-Range,Content-Disposition")

		next.ServeHTTP(w, r)
	})
}
