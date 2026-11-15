package middleware

import (
	"context"
	"errors"
	"m-drobynin/go-ext-url-shortener/internal/model"
	"m-drobynin/go-ext-url-shortener/internal/service"
	"net/http"

	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

const authCookieName = "Authorization"
const secretKey = "supersecretkey"

var ErrInvalidJWT = errors.New("invalid jwt")
var ErrUnexpectedSigning = errors.New("unexpected signing")
var ErrUnauthorized = errors.New("unauthorized")

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewJwtAuthorizationMiddleware(usersService service.UsersService) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(NewJwtHandler(h, usersService))
	}
}

func NewJwtHandler(h http.Handler, usersService service.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := processJWTAuth(r, usersService)

		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		token, err := buildJWT(user.UserID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  authCookieName,
			Value: token,
		})
		ctx := context.WithValue(r.Context(), model.UserIDKey, user.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func processJWTAuth(r *http.Request, usersService service.UsersService) (*model.User, error) {
	cookie, err := r.Cookie(authCookieName)

	if err != nil {
		return usersService.Create(r.Context())
	}

	claims, err := getTokenClaims(cookie.Value)

	if err != nil {
		if errors.Is(err, ErrInvalidJWT) {
			return usersService.Create(r.Context())
		}

		if errors.Is(err, ErrUnexpectedSigning) {
			return nil, ErrUnauthorized
		}

		return nil, err
	}

	if claims.UserID == "" {
		return nil, ErrUnauthorized
	}

	user, err := usersService.GetByID(r.Context(), claims.UserID)

	if err != nil {
		if errors.Is(err, model.ErrDBNotFound) {
			return nil, ErrUnauthorized
		}

		return nil, err
	}

	return user, nil
}

func buildJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigning, t.Header["alg"])
			}

			return []byte(secretKey), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
