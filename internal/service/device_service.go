package service

import (
	"time"

	"github.com/fleet-monitor/mini-fleet/internal/model"
	"github.com/fleet-monitor/mini-fleet/internal/repository"
)

type DeviceService struct {
	repo           repository.DeviceRepository
	timeoutDuration time.Duration
}

func NewDeviceService(repo repository.DeviceRepository, timeout time.Duration) *DeviceService {
	return &DeviceService{
		repo:           repo,
		timeoutDuration: timeout,
	}
}

func (s *DeviceService) RegisterDevice(req model.RegisterDeviceRequest) (*model.Device, error) {
	if req.ID == "" || req.Name == "" {
		return nil, repository.ErrDeviceNotFound
	}

	dev := &model.Device{
		ID:           req.ID,
		Name:         req.Name,
		Status:       model.StatusOffline,
		RegisteredAt: time.Now().UTC(),
	}

	if err := s.repo.Save(dev); err != nil {
		return nil, err
	}

	return dev, nil
}

func (s *DeviceService) RecordHeartbeat(id string, req model.HeartbeatRequest) (*model.Device, error) {
	hbTime := req.Timestamp
	if hbTime.IsZero() {
		hbTime = time.Now().UTC()
	}

	dev, err := s.repo.UpdateHeartbeat(id, hbTime, req.CPUUsage, req.SignalStrength)
	if err != nil {
		return nil, err
	}

	s.computeStatus(dev)
	return dev, nil
}

func (s *DeviceService) GetDevice(id string) (*model.Device, error) {
	dev, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.computeStatus(dev)
	return dev, nil
}

func (s *DeviceService) ListDevices() []*model.Device {
	devices := s.repo.FindAll()
	for _, dev := range devices {
		s.computeStatus(dev)
	}
	return devices
}

func (s *DeviceService) GetFleetSummary() model.FleetSummary {
	devices := s.ListDevices()
	summary := model.FleetSummary{
		Total: len(devices),
	}

	for _, dev := range devices {
		if dev.Status == model.StatusOnline {
			summary.Online++
		} else {
			summary.Offline++
		}
	}

	return summary
}

// Dynamic evaluation rule: heartbeat within last 30 seconds -> ONLINE, else OFFLINE
func (s *DeviceService) computeStatus(dev *model.Device) {
	if dev.LastHeartbeat == nil {
		dev.Status = model.StatusOffline
		return
	}

	if time.Since(*dev.LastHeartbeat) <= s.timeoutDuration {
		dev.Status = model.StatusOnline
	} else {
		dev.Status = model.StatusOffline
	}
}