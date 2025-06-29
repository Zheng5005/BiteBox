package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(r *http.Request, secret string) (string, error)  {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("Mssing token")
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id missing")
	}

	return userID, nil
}

func GenerateMockJWT(userID, secret string) (string, error)  {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}
