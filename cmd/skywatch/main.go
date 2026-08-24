package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1260124186-cc/skywatch-coordinator/internal/httpapi"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func main() {
	addr := os.Getenv("SKYWATCH_ADDR")
	if addr == "" {
		addr = "127.0.0.1:18087"
	}
	repository := store.NewMemoryStore()
	coordinator := service.NewCoordinator(repository)
	if err := service.LoadDemo(coordinator); err != nil {
		log.Fatal(err)
	}
	server := httpapi.NewServer(coordinator)
	log.Printf("skywatch coordinator listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
