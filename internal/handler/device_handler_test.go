package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fleet-monitor/mini-fleet/internal/model"
	"github.com/fleet-monitor/mini-fleet/internal/repository"
	"github.com/fleet-monitor/mini-fleet/internal/service"
)

func setupTestServer() *DeviceHandler {
	repo := repository.NewInMemoryDeviceRepository()
	svc := service.NewDeviceService(repo, 30*time.Second)
	return NewDeviceHandler(svc)
}

func TestRegisterDeviceHandler(t *testing.T) {
	h := setupTestServer()

	payload := []byte(`{"id":"dev-101","name":"Test Device"}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	h.RegisterDevice(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", resp.StatusCode)
	}

	var dev model.Device
	_ = json.NewDecoder(resp.Body).Decode(&dev)

	if dev.ID != "dev-101" || dev.Name != "Test Device" {
		t.Fatalf("unexpected device data returned: %+v", dev)
	}
}

func TestGetSummaryHandler(t *testing.T) {
	h := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	w := httptest.NewRecorder()

	h.GetSummary(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}
}