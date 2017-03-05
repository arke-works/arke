package middleware

import (
	"context"
	"iris.arke.works/forum/http/ctxkeys"
	"iris.arke.works/forum/http/helper"
	"net/http"
	"strconv"
)

func PageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			pivotID, size int64
			err           error
		)

		pivotString := r.URL.Query().Get("pivot_id")
		sizeString := r.URL.Query().Get("req_size")

		if pivotString != "" {
			pivotID, err = strconv.ParseInt(pivotString, 10, 63)
			if err != nil {
				helper.ErrorWriter(w, r, http.StatusBadRequest, err)
				return
			}
		}

		if sizeString != "" {
			size, err = strconv.ParseInt(sizeString, 10, 63)
			if err != nil {
				helper.ErrorWriter(w, r, http.StatusBadRequest, err)
				return
			}
		}

		if size > 25 {
			size = 25
		}

		r = r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxPivotIDKey, pivotID))
		r = r.WithContext(context.WithValue(r.Context(), ctxkeys.CtxSizeKey, size))

		next.ServeHTTP(w, r)
	})
}
