package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/techops-interviews/service-registry/internal/api"
	"github.com/techops-interviews/service-registry/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.New()
	h := api.NewHandler(s)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting service registry on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
