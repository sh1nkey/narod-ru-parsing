package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			dur := fmt.Sprintf("%dms", time.Duration(time.Since(start).Milliseconds()))

			log.Trace().
				Str("method", r.Method).
				Str("host", r.Host).
				Str("uri", r.RequestURI).
				Str("proto", r.Proto).
				Str("origin", r.Header.Get("Origin")).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Str("duration", dur).Send()
		}()
		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}