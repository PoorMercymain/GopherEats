package token

import "github.com/golang-jwt/jwt/v4"

func JWT(email string, passwordHash string, jwtKey string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"hash":  passwordHash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
