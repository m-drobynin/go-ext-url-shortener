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

var ErrInvalidJWT = errors.New("invalid jwt")
var ErrUnexpectedSigning = errors.New("unexpected signing")
var ErrUnauthorized = errors.New("unauthorized")

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewJwtAuthorizationMiddleware(usersService service.UsersService, jwtSecret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(NewJwtHandler(h, usersService, jwtSecret))
	}
}

func NewJwtHandler(h http.Handler, usersService service.UsersService, jwtSecret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := processJWTAuth(r, usersService, jwtSecret)

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

		token, err := buildJWT(user.UserID, jwtSecret)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		ctx := context.WithValue(r.Context(), model.UserIDKey, user.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func processJWTAuth(r *http.Request, usersService service.UsersService, jwtSecret string) (*model.User, error) {
	cookie, err := r.Cookie(authCookieName)

	if err != nil {
		return usersService.Create(r.Context())
	}

	claims, err := getTokenClaims(cookie.Value, jwtSecret)

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

func buildJWT(userID string, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenClaims(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigning, t.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidJWT, err)
	}

	if !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
