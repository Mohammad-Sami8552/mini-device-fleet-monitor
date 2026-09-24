package service

import (
	"testing"
	"time"

	"github.com/fleet-monitor/mini-fleet/internal/model"
	"github.com/fleet-monitor/mini-fleet/internal/repository"
)

func TestTimeoutLogic(t *testing.T) {
	repo := repository.NewInMemoryDeviceRepository()
	timeout := 2 * time.Second // Short timeout for unit testing
	svc := NewDeviceService(repo, timeout)

	// 1. Register device
	_, _ = svc.RegisterDevice(model.RegisterDeviceRequest{ID: "d1", Name: "Device 1"})

	// Initially offline
	d, _ := svc.GetDevice("d1")
	if d.Status != model.StatusOffline {
		t.Fatalf("expected initial status OFFLINE, got %s", d.Status)
	}

	// 2. Heartbeat sent -> ONLINE
	now := time.Now()
	_, _ = svc.RecordHeartbeat("d1", model.HeartbeatRequest{Timestamp: now})
	d, _ = svc.GetDevice("d1")
	if d.Status != model.StatusOnline {
		t.Fatalf("expected status ONLINE, got %s", d.Status)
	}

	// 3. Wait past timeout -> OFFLINE
	time.Sleep(2100 * time.Millisecond)
	d, _ = svc.GetDevice("d1")
	if d.Status != model.StatusOffline {
		t.Fatalf("expected status OFFLINE after timeout, got %s", d.Status)
	}
}

func TestFleetSummary(t *testing.T) {
	repo := repository.NewInMemoryDeviceRepository()
	svc := NewDeviceService(repo, 30*time.Second)

	svc.RegisterDevice(model.RegisterDeviceRequest{ID: "d1", Name: "D1"})
	svc.RegisterDevice(model.RegisterDeviceRequest{ID: "d2", Name: "D2"})

	svc.RecordHeartbeat("d1", model.HeartbeatRequest{Timestamp: time.Now()})

	summary := svc.GetFleetSummary()
	if summary.Total != 2 || summary.Online != 1 || summary.Offline != 1 {
		t.Fatalf("incorrect summary: %+v", summary)
	}
}