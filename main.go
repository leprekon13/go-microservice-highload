package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go-microservice-highload/handlers"
	"go-microservice-highload/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userSvc := services.NewUserService()
	userHandler := handlers.NewUserHandler(userSvc)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", userHandler.GetUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/users", userHandler.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
