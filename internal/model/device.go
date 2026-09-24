package model

import "time"

type DeviceStatus string

const (
	StatusOnline  DeviceStatus = "ONLINE"
	StatusOffline DeviceStatus = "OFFLINE"
)

// RegisterDeviceRequest holds payload for new device registration
type RegisterDeviceRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HeartbeatRequest holds payload sent by active devices
type HeartbeatRequest struct {
	Timestamp     time.Time `json:"timestamp"`
	Status        string    `json:"status"`
	CPUUsage      *float64  `json:"cpu_usage,omitempty"`
	SignalStrength *int     `json:"signal_strength,omitempty"`
}

// Device represents the complete internal and external device state
type Device struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Status        DeviceStatus `json:"status"`
	LastHeartbeat *time.Time   `json:"last_heartbeat"`
	CPUUsage      *float64     `json:"cpu_usage,omitempty"`
	SignalStrength *int        `json:"signal_strength,omitempty"`
	RegisteredAt  time.Time    `json:"registered_at"`
}

// FleetSummary holds overall health metrics of the fleet
type FleetSummary struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}