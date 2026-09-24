package repository

import (
	"testing"
	"time"

	"github.com/fleet-monitor/mini-fleet/internal/model"
)

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryDeviceRepository()

	// Test Registration
	dev := &model.Device{ID: "d1", Name: "Device 1", RegisteredAt: time.Now()}
	err := repo.Save(dev)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test Duplicate Registration
	err = repo.Save(dev)
	if err != ErrDeviceAlreadyExists {
		t.Fatalf("expected ErrDeviceAlreadyExists, got %v", err)
	}

	// Test Get By ID
	found, err := repo.FindByID("d1")
	if err != nil || found.Name != "Device 1" {
		t.Fatalf("failed to retrieve device properly")
	}

	// Test Heartbeat Update
	now := time.Now()
	cpu := 45.5
	updated, err := repo.UpdateHeartbeat("d1", now, &cpu, nil)
	if err != nil || updated.LastHeartbeat == nil || *updated.CPUUsage != 45.5 {
		t.Fatalf("failed to update heartbeat")
	}
}