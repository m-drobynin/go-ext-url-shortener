package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

func main() {
	tokenString, err := BuildJWTString()

	if err != nil {
		log.Fatal(err)
	}

	userId := GetUserID(tokenString)

	fmt.Println(userId)
	fmt.Println(tokenString)
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		// когда создан токен
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		// собственное утверждение
		UserID: 1,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) int {
	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token is valid")

	// возвращаем ID пользователя в читаемом виде
	return claims.UserID
}
