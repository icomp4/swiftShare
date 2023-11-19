package validators

import (
	"fmt"
	"net/http"
	"os"
	"unicode"
	"github.com/golang-jwt/jwt/v5"
)

func PasswordIsValid(password string) bool {
	if len(password) < 6 {
		return false
	}
	var (
		hasUpper, hasLower, hasDigit, hasSpecial bool
	)
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if isSpecial(char) {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func isSpecial(char rune) bool {
	specialCharacters := "!@#$%^&*()_+{}[]|:;<>,.?/~"
	for _, special := range specialCharacters {
		if char == special {
			return true
		}
	}
	return false
}

func ExtractToken(r *http.Request) error {
	tokenCookie, err := r.Cookie("jwt")
	if err != nil {
		return err
	}
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
	}
	return nil
}
