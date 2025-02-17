package api

import (
	"context"
	"net/http"
	"strconv"
)

const DefaultPageLimit = 50

type CtxKey string

const (
	CtxKeyLimit       CtxKey = "limit"
	CtxKeyOffset      CtxKey = "offset"
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		var err error
		limit := DefaultPageLimit
		if limitStr !=  "" {
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




