package gapi

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func HttpRecovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Any("panic", r).
					Bytes("stack", debug.Stack()).
					Msg("recovered from panic")

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
