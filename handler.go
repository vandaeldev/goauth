package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	result := GetDB().Find(&users)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(result.Error.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	result := GetDB().First(&user, "id = ?", id)
	if result.Error != nil {
		HandleResultErr(w, result)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	hashed, err := HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user.Password = hashed
	result := GetDB().Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(result.Error.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		creds User
		user  User
	)
	json.NewDecoder(r.Body).Decode(&creds)
	result := GetDB().First(&user, "email = ?", creds.Email)
	if result.Error != nil {
		HandleResultErr(w, result)
		return
	}
	if CheckPassword(user.Password, creds.Password) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	t, err := GenerateJWT(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res := fmt.Sprintf("{\"token\": \"%s\"}", t)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	userResult := GetDB().First(&user, "id = ?", id)
	if userResult.Error != nil {
		HandleResultErr(w, userResult)
		return
	}
	var updatedUser User
	json.NewDecoder(r.Body).Decode(&updatedUser)
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	updatedUser.ID = uint(parsedID)
	result := GetDB().Save(&updatedUser)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(result.Error.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result := GetDB().Delete(&User{}, id)
	if result.Error != nil {
		HandleResultErr(w, result)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
