package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"rapidEx/internal/domain/user"
)
//TODO: add dynamic salt
func NewToken(user user.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.UUID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration)

	tokenString, err := token.SignedString([]byte("A"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}