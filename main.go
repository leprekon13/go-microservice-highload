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
	"go-microservice-highload/metrics"
	"go-microservice-highload/services"
	"go-microservice-highload/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint == "" {
		s3Endpoint = "localhost:9000"
	}
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	if s3AccessKey == "" {
		s3AccessKey = "minioadmin"
	}
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	if s3SecretKey == "" {
		s3SecretKey = "minioadmin"
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "users-data"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	async := utils.NewAsyncProcessor(10000, 10000)
	async.Start(ctx)

	limiter := utils.NewRateLimiter(1000, 5000)

	userSvc := services.NewUserService()

	intSvc, err := services.NewIntegrationService(s3Endpoint, s3AccessKey, s3SecretKey, s3Bucket)
	if err != nil {
		log.Fatalf("failed to connect s3: %v", err)
	}
	if err := intSvc.InitBucket(ctx); err != nil {
		log.Printf("warning: could not init s3 bucket: %v", err)
	}

	userHandler := handlers.NewUserHandler(userSvc, async)
	intHandler := handlers.NewIntegrationHandler(userSvc, intSvc)

	r := mux.NewRouter()
	r.Handle("/metrics", metrics.Handler()).Methods(http.MethodGet)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(limiter.Middleware)
	api.Use(metrics.Middleware)

	api.HandleFunc("/users", userHandler.GetUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/users", userHandler.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)
	api.HandleFunc("/integration/export", intHandler.Export).Methods(http.MethodPost)

	addr := ":" + port
	log.Printf("listening on %s", addr)

	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	server.Shutdown(context.Background())
}
