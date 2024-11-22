package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	defer CloseDB()
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/users", GetUsers).Methods("GET")
	s.HandleFunc("/users/{id}", GetUser).Methods("GET")
	s.HandleFunc("/users", CreateUser).Methods("POST")
	s.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	s.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8888", r))
}
