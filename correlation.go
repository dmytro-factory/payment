package payment

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
)

type correlationIDKey struct{}

func getCorrelationID(r *http.Request) string {
	cid := r.Header.Get("X-Correlation-ID")
	if cid == "" {
		b := make([]byte, 16)
		rand.Read(b)
		cid = fmt.Sprintf("%x", b)
	}
	return cid
}

func correlationIDMiddleware(logger log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid := getCorrelationID(r)
		w.Header().Set("X-Correlation-ID", cid)
		ctx := context.WithValue(r.Context(), correlationIDKey{}, cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
