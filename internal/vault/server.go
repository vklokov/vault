package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RecordParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Server() {
	config := NewConfig()
	vault := New(config)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		records := vault.All()

		w.WriteHeader(200)

		renderOK(w, records)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var params RecordParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			log.Fatalf("failed to encode params: %v", err)
			render400(w)
			return
		}
		defer r.Body.Close()

		err := vault.Upsert(params.Key, params.Value)
		if err != nil {
			log.Fatalf("failed to upsert a record: %v", err)
			render500(w)
			return
		}

		w.WriteHeader(200)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 server started on port: %s", config.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting the server: %v", err)
		}
	}()

	<-stop

	log.Println("🚦 shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down the server: %v", err)
	}

	log.Println("🛑 server stopped")
}

func respondWith(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}

func renderOK(w http.ResponseWriter, payload interface{}) {
	respondWith(w, http.StatusOK, payload)
}

func render400(w http.ResponseWriter) {
	respondWith(w, http.StatusBadRequest, map[string]string{
		"error": "Bad request",
	})
}

func render500(w http.ResponseWriter) {
	respondWith(w, http.StatusBadRequest, map[string]string{
		"error": "Internal Server Error",
	})
}
