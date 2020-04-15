package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken -> create JWT token
func CreateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))

}

// TokenValidation ...
func TokenValidation(bearerToken string) (uint, error) {
	tokenString := strings.Split(bearerToken, " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		var uid uint = uint(claims["userId"].(float64))
		return uid, nil
	}

	return 0, nil
}
