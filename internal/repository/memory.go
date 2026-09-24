package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/fleet-monitor/mini-fleet/internal/model"
)

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already registered")
)

type DeviceRepository interface {
	Save(device *model.Device) error
	FindByID(id string) (*model.Device, error)
	FindAll() []*model.Device
	UpdateHeartbeat(id string, timestamp time.Time, cpu *float64, signal *int) (*model.Device, error)
}

type InMemoryDeviceRepository struct {
	mu      sync.RWMutex
	devices map[string]*model.Device
}

func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*model.Device),
	}
}

func (r *InMemoryDeviceRepository) Save(device *model.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return ErrDeviceAlreadyExists
	}

	r.devices[device.ID] = device
	return nil
}

func (r *InMemoryDeviceRepository) FindByID(id string) (*model.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dev, exists := r.devices[id]
	if !exists {
		return nil, ErrDeviceNotFound
	}

	// Return a copy to ensure thread-safety on returned object
	copied := *dev
	return &copied, nil
}

func (r *InMemoryDeviceRepository) FindAll() []*model.Device {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*model.Device, 0, len(r.devices))
	for _, dev := range r.devices {
		copied := *dev
		list = append(list, &copied)
	}
	return list
}

func (r *InMemoryDeviceRepository) UpdateHeartbeat(id string, timestamp time.Time, cpu *float64, signal *int) (*model.Device, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dev, exists := r.devices[id]
	if !exists {
		return nil, ErrDeviceNotFound
	}

	dev.LastHeartbeat = &timestamp
	if cpu != nil {
		dev.CPUUsage = cpu
	}
	if signal != nil {
		dev.SignalStrength = signal
	}

	copied := *dev
	return &copied, nil
}