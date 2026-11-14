package middleware

import (
	"net/http"
)

func WithJwtAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("Authorization")

		// if err != nil && errors.Is(err, http.ErrNoCookie) {

		// }

		h.ServeHTTP(w, r)
	})
}
