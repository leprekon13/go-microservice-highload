package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"go-microservice-highload/handlers"
	"go-microservice-highload/services"
	"go-microservice-highload/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	async := utils.NewAsyncProcessor(10000, 10000)
	async.Start(ctx)

	userSvc := services.NewUserService()
	userHandler := handlers.NewUserHandler(userSvc, async)

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
