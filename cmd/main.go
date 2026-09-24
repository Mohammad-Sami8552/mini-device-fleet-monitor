package main

import (
	"log"
	"net/http"

	"github.com/fleet-monitor/mini-fleet/internal/config"
	"github.com/fleet-monitor/mini-fleet/internal/handler"
	"github.com/fleet-monitor/mini-fleet/internal/repository"
	"github.com/fleet-monitor/mini-fleet/internal/service"
)

func main() {
	cfg := config.Load()

	repo := repository.NewInMemoryDeviceRepository()
	svc := service.NewDeviceService(repo, cfg.TimeoutDuration)
	h := handler.NewDeviceHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /devices", h.RegisterDevice)
	mux.HandleFunc("GET /devices", h.ListDevices)
	mux.HandleFunc("GET /summary", h.GetSummary)
	
	// Handles both GET /devices/{id} and POST /devices/{id}/heartbeat
	mux.HandleFunc("/devices/", h.HandleDeviceByID)

	log.Printf("Starting Fleet Monitor Server on port :%s ...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}