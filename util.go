package main

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func HandleResultErr(w http.ResponseWriter, result *gorm.DB) {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(result.Error.Error()))
}
