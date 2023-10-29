package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/PoorMercymain/GopherEats/internal/app/auth/errors"
)

func JWT(email string, passwordHash string, jwtKey string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"hash":  passwordHash,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetAuthDataFromJWT(token string, jwtKey string) (string, string, error) {
	claims := jwt.MapClaims{
		"email": "",
		"hash":  "",
		"exp":   time.Now().Unix(),
	}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", "", errors.ErrorWrongToken
	}

	if !jwtToken.Valid {
		return "", "", errors.ErrorWrongToken
	}

	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		return "", "", errors.ErrorWrongToken
	}

	return claims["email"].(string), claims["hash"].(string), nil
}
