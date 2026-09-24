package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fleet-monitor/mini-fleet/internal/model"
	"github.com/fleet-monitor/mini-fleet/internal/repository"
	"github.com/fleet-monitor/mini-fleet/internal/service"
)

type DeviceHandler struct {
	service *service.DeviceService
}

func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: svc}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func (h *DeviceHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.ID) == "" || strings.TrimSpace(req.Name) == "" {
		respondError(w, http.StatusBadRequest, "id and name are required")
		return
	}

	dev, err := h.service.RegisterDevice(req)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceAlreadyExists) {
			respondError(w, http.StatusConflict, "device already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, dev)
}

func (h *DeviceHandler) HandleDeviceByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/devices/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		respondError(w, http.StatusBadRequest, "device id is required")
		return
	}

	id := parts[0]

	// Handle POST /devices/{id}/heartbeat
	if len(parts) == 2 && parts[1] == "heartbeat" {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var hb model.HeartbeatRequest
		if err := json.NewDecoder(r.Body).Decode(&hb); err != nil {
			respondError(w, http.StatusBadRequest, "invalid heartbeat payload")
			return
		}

		dev, err := h.service.RecordHeartbeat(id, hb)
		if err != nil {
			if errors.Is(err, repository.ErrDeviceNotFound) {
				respondError(w, http.StatusNotFound, "device not registered")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, dev)
		return
	}

	// Handle GET /devices/{id}
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		dev, err := h.service.GetDevice(id)
		if err != nil {
			if errors.Is(err, repository.ErrDeviceNotFound) {
				respondError(w, http.StatusNotFound, "device not found")
				return
			}
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusOK, dev)
		return
	}

	respondError(w, http.StatusNotFound, "endpoint not found")
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	devices := h.service.ListDevices()
	respondJSON(w, http.StatusOK, devices)
}

func (h *DeviceHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	summary := h.service.GetFleetSummary()
	respondJSON(w, http.StatusOK, summary)
}