package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type DeviceSim struct {
	ID    string
	Name  string
	Delay time.Duration
}

func main() {
	// Agar aapka server port 8081 par chal raha hai toh baseURL ko "http://localhost:8081" kar do
	baseURL := "http://localhost:8081"
	client := &http.Client{Timeout: 5 * time.Second}

	devices := []DeviceSim{
		{ID: "device-01", Name: "Lab Device 01", Delay: 5 * time.Second},
		{ID: "device-02", Name: "Lab Device 02", Delay: 5 * time.Second},
		{ID: "device-03", Name: "Lab Device 03", Delay: 5 * time.Second},
		{ID: "device-04", Name: "Lab Device 04", Delay: 5 * time.Second},
		{ID: "device-05", Name: "Lab Device 05 (Will Stop Early)", Delay: 5 * time.Second},
	}

	log.Println("--- Starting Fleet Simulator ---")

	// 1. Register all devices
	for _, dev := range devices {
		payload, _ := json.Marshal(map[string]string{
			"id":   dev.ID,
			"name": dev.Name,
		})

		resp, err := client.Post(baseURL+"/devices", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalf("Failed to register %s: %v", dev.ID, err)
		}
		resp.Body.Close()
		log.Printf("Registered device: %s (%s)", dev.ID, dev.Name)
	}

	var wg sync.WaitGroup

	// 2. Start active heartbeat routines
	for idx, dev := range devices {
		wg.Add(1)
		go func(d DeviceSim, index int) {
			defer wg.Done()

			ticker := time.NewTicker(d.Delay)
			defer ticker.Stop()

			count := 0
			for range ticker.C {
				count++

				// Simulate failure for device-05 after 3 heartbeats (15 seconds)
				if d.ID == "device-05" && count > 3 {
					log.Printf("⚠️ STOPPING HEARTBEATS FOR %s (Simulating failure / offline status)...", d.ID)
					return
				}

				cpu := 20.0 + rand.Float64()*30.0
				sig := -80 + rand.Intn(30)

				hbPayload, _ := json.Marshal(map[string]interface{}{
					"timestamp":       time.Now().UTC().Format(time.RFC3339),
					"status":          "OK",
					"cpu_usage":       cpu,
					"signal_strength": sig,
				})

				url := fmt.Sprintf("%s/devices/%s/heartbeat", baseURL, d.ID)
				resp, err := client.Post(url, "application/json", bytes.NewBuffer(hbPayload))
				if err != nil {
					log.Printf("Error sending heartbeat for %s: %v", d.ID, err)
					continue
				}
				resp.Body.Close()

				log.Printf("Sent heartbeat #%d for %s [CPU: %.1f%%, Signal: %ddBm]", count, d.ID, cpu, sig)
			}
		}(dev, idx)
	}

	wg.Wait()
	log.Println("Simulator finished routine. Keep server running to inspect /summary endpoint.")
}