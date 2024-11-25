package main

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

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
