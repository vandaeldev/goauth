package main

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateKeyPair() {
	dirname := "keys"
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirname, "pub.pem"), pub, 0644); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirname, "priv.pem"), priv, 0640); err != nil {
		log.Fatal(err)
	}
}

func HandleResultErr(w http.ResponseWriter, result *gorm.DB) {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(result.Error.Error()))
}

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	return string(bytes), err
}

func CheckPassword(hash string, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err != nil
}

func GenerateJWT(user User) (string, error) {
	key, err := os.ReadFile(filepath.Join("keys", "priv.pem"))
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": "dev.vandael/goauth",
		"sub": user.Email,
		"id":  user.ID,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(time.Minute * 30).Unix(),
	}
	t := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)
	return t.SignedString(ed25519.PrivateKey(key))
}

func ValidateJWT(token string) (*jwt.Token, error) {
	return jwt.NewParser(
		jwt.WithIssuer("dev.vandael/goauth"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{"EdDSA"}),
	).Parse(token, withKey)
}

func withKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	key, err := os.ReadFile(filepath.Join("keys", "pub.pem"))
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(key), nil
}
