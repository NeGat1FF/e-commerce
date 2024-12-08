package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateJWT(claims map[string]interface{}, secretKey string) (string, error) {
	jwtClaims := jwt.MapClaims{}

	for key, value := range claims {
		jwtClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func ValidateClaims(claims jwt.MapClaims, tokenType string) bool {
	if claims["type"] != tokenType {
		return false
	}

	if claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return false
	}

	if _, ok := claims["uid"]; !ok {
		return false
	}

	return uuid.Validate(claims["uid"].(string)) == nil
}
