package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()
	defer CloseDB()
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/users", GetUsers).Methods("GET")
	s.HandleFunc("/users/{id}", GetUser).Methods("GET")
	s.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	s.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	s.HandleFunc("/signup", Signup).Methods("POST")
	s.HandleFunc("/login", Login).Methods("POST")
	srv := &http.Server{
		Addr:         "0.0.0.0:8888",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down...")
	os.Exit(0)
}
