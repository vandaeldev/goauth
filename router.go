package main

import "github.com/gorilla/mux"

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter := apiRouter.PathPrefix("/users").Subrouter()
	apiRouter.HandleFunc("/signup", Signup).Methods("POST")
	apiRouter.HandleFunc("/login", Login).Methods("POST")
	subRouter.HandleFunc("", GetUsers).Methods("GET")
	subRouter.HandleFunc("/{id}", GetUser).Methods("GET")
	subRouter.HandleFunc("/{id}", UpdateUser).Methods("PUT")
	subRouter.HandleFunc("/{id}", DeleteUser).Methods("DELETE")
	subRouter.Use(VerifyToken)
	return router
}
