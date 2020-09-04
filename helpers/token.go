package helpers

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ibamibrhm/donation-server/models"
)

// PayloadToken ...
type PayloadToken struct {
	UserID       uint
	TokenVersion uint
}

// CreateToken -> create JWT token
func CreateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = user.ID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() //Token expires after 15 minutes
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// CreateRefreshToken -> create refresh JWT token
func CreateRefreshToken(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = user.ID
	claims["tokenVersion"] = user.TokenVersion
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() //Token expires after 7 days
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

// TokenValidation -> validate JWT token
func TokenValidation(bearerToken string, secret string) (PayloadToken, error) {
	tokenString := strings.Split(bearerToken, " ")

	if len(tokenString) != 2 {
		return PayloadToken{}, errors.New("Invalid token")
	}

	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return PayloadToken{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		var payload PayloadToken

		payload.UserID = uint(claims["userId"].(float64))

		if val, ok := claims["tokenVersion"]; ok {
			payload.TokenVersion = uint(val.(float64))
		}

		return payload, nil
	}

	return PayloadToken{}, errors.New("Unexpected error")
}
