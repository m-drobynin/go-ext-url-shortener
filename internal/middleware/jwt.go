package middleware

import (
	"context"
	"net/http"
)

func WithJwtAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("Authorization")

		// if err != nil && errors.Is(err, http.ErrNoCookie) {

		// }

		ctx := context.WithValue(r.Context(), "user_id", "user_id")
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	})
}
