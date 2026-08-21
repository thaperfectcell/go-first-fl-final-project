package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 8 * time.Hour

type signinRequest struct {
	Password string `json:"password"`
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func CreateToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"pwd": passwordHash(password),
		"exp": time.Now().Add(tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("todo-auth:" + password)
	return token.SignedString(secret)
}

func validateToken(tokenStr, password string) bool {
	if tokenStr == "" {
		return false
	}

	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("todo-auth:" + password), nil
	})
	if err != nil || !parsed.Valid {
		return false
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	storedHash, ok := claims["pwd"].(string)
	if !ok {
		return false
	}

	return storedHash == passwordHash(password)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := todoPassword
		if pass == "" {
			next(w, r)
			return
		}

		jwtCookie, err := r.Cookie("token")
		if err != nil || !validateToken(jwtCookie.Value, pass) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
