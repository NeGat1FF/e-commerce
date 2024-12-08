package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// ValidateToken validates the token returning claims if the token is valid
func ValidateToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("error validating token: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	if ok && token.Valid && claims.VerifyExpiresAt(jwt.TimeFunc().Unix(), true) && claims["type"].(string) == "access" {
		return claims, nil
	}

	return nil, err
}
