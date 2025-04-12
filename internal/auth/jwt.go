package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Claims struct {
	jwt.RegisteredClaims
	Sub        string `json:"sub"`
	Department string `json:"department"`
	Role       string `json:"role"`
}

const USER_ID = "uid"
const NAME = "department"

func GenerateJWT(sub string, dept string, rol string) (token string, err error) {
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		sub,
		dept,
		rol,
	}
	token, err = createToken(claims)
	return token, errors.Wrap(err, "fail to create token")
}

func ParseToken(token string) (Claims, error) {
	claims := Claims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if secret, err := model.GetSettings(context.TODO(), "JWTSecret"); err == nil {
			return []byte(secret), nil
		} else {
			return nil, fmt.Errorf("fail to get JWT secret: %v", err)
		}
	})
	return claims, errors.Wrap(err, "fail to parse token")
}

func createToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := model.GetSettings(context.TODO(), "JWTSecret")
	if err != nil {
		return "", fmt.Errorf("fail to get JWT secret: %v", err)
	}
	return token.SignedString([]byte(secret))
}
