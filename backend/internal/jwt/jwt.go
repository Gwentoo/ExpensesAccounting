package jwt

import (
	"backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(id int64, email string, duration time.Duration) (string, error) {
	cfg := config.LoadConfig()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = id
	claims["email"] = email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
