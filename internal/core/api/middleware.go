package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

const DefaultPageLimit = 50

type CtxKey string

const (
	CtxKeyLimit  CtxKey = "limit"
	CtxKeyOffset CtxKey = "offset"
	CtxKeyUser   CtxKey = "user"
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		var err error
		limit := DefaultPageLimit
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				limit = DefaultPageLimit
			}
		}

		offset := 0
		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				offset = 0
			}
		}

		ctx := context.WithValue(r.Context(), CtxKeyLimit, limit)
		ctx = context.WithValue(ctx, CtxKeyOffset, offset)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
