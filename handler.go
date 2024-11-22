package main

import (
	"encoding/json"
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
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	result := GetDB().First(&user, "id = ?", id)
	if result.Error != nil {
		HandleResultErr(w, result)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	result := GetDB().Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(result.Error.Error()))
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user User
	result := GetDB().First(&user, "id = ?", id)
	if result.Error != nil {
		HandleResultErr(w, result)
	} else {
		var updatedUser User
		json.NewDecoder(r.Body).Decode(&updatedUser)
		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			updatedUser.ID = uint(parsedID)
			result := GetDB().Save(&updatedUser)
			if result.Error != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(result.Error.Error()))
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedUser)
			}
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result := GetDB().Delete(&User{}, id)
	if result.Error != nil {
		HandleResultErr(w, result)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
