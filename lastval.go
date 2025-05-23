package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	PROXMOX_HOST       = "192.168.1.2"
	PROXMOX_PORT       = "8006"
	API_USER           = "rdn@pve"
	API_TOKEN_NAME     = "n7h6HPHJHu7sx8bGtbEcdnyW"
	API_TOKEN_SECRET   = "ccc3bccc-5466-465e-a4b0-68e9f553bb87"
	TARGET_VMID        = 201
	POLL_INTERVAL_SEC  = 1
	DATA_DIR           = "data"
	BANDWIDTH_INI_FILE = "LastBandwidth.ini"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	lastNetIn      float64
	lastNetOut     float64
	lastPercentage float64
	prevNetIn      float64
	prevNetOut     float64
	prevTime       time.Time
)

func ensureDataDir() error {
	if _, err := os.Stat(DATA_DIR); os.IsNotExist(err) {
		return os.Mkdir(DATA_DIR, 0755)
	}
	return nil
}

func getIniPath() string {
	return filepath.Join(DATA_DIR, BANDWIDTH_INI_FILE)
}

func loadLastValues() error {
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	path := getIniPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Initialize with empty values
		lastNetIn = 0
		lastNetOut = 0
		lastPercentage = 0
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read ini file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "netIn=") {
			if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "netIn="), 64); err == nil {
				lastNetIn = val
			}
		}
		if strings.HasPrefix(line, "netOut=") {
			if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "netOut="), 64); err == nil {
				lastNetOut = val
			}
		}
		if strings.HasPrefix(line, "percentage=") {
			if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "percentage="), 64); err == nil {
				lastPercentage = val
			}
		}
	}
	return nil
}

func saveCurrentValues(netIn, netOut, percentage float64) {
	path := getIniPath()
	data := fmt.Sprintf("netIn=%.2f\nnetOut=%.2f\npercentage=%.2f\n", netIn, netOut, percentage)
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		fmt.Printf("Warning: failed to save values: %v\n", err)
	}
}

func getVMStats() (float64, float64, error) {
	url := fmt.Sprintf("https://%s:%s/api2/json/cluster/resources?type=vm", PROXMOX_HOST, PROXMOX_PORT)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("PVEAPIToken=%s!%s=%s",
		API_USER, API_TOKEN_NAME, API_TOKEN_SECRET))

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("API error: %s", resp.Status)
	}

	var result struct {
		Data []struct {
			Type   string  `json:"type"`
			VMID   int     `json:"vmid"`
			NetIn  float64 `json:"netin"`
			NetOut float64 `json:"netout"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	for _, resource := range result.Data {
		if resource.Type == "qemu" && resource.VMID == TARGET_VMID {
			return resource.NetIn, resource.NetOut, nil
		}
	}

	return 0, 0, fmt.Errorf("VM %d not found", TARGET_VMID)
}

func main() {
	// Load previous values
	if err := loadLastValues(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	fmt.Println("Time\t\tRX (B/s)\tTX (B/s)\tUsage (%)")

	// Initial data
	var err error
	prevNetIn, prevNetOut, err = getVMStats()
	if err != nil {
		fmt.Printf("Initial data error: %v\n", err)
		os.Exit(1)
	}
	prevTime = time.Now()

	// Wait for initial non-zero values
	for prevNetIn == 0 && prevNetOut == 0 {
		time.Sleep(POLL_INTERVAL_SEC * time.Second)
		prevNetIn, prevNetOut, err = getVMStats()
		if err != nil {
			fmt.Printf("Error getting initial non-zero values: %v\n", err)
			continue
		}
	}

	for {
		time.Sleep(POLL_INTERVAL_SEC * time.Second)

		currentNetIn, currentNetOut, err := getVMStats()
		if err != nil {
			fmt.Printf("Error: %v - Using last known values\n", err)
			fmt.Printf("%s\t%.2f\t%.2f\t%.2f%%\n",
				time.Now().Format("15:04:05"),
				lastNetIn,
				lastNetOut,
				lastPercentage)
			continue
		}

		// Skip if we got zero values
		if currentNetIn == 0 && currentNetOut == 0 {
			continue
		}

		deltaTime := time.Since(prevTime).Seconds()
		rxBps := (currentNetIn - prevNetIn) / deltaTime
		txBps := (currentNetOut - prevNetOut) / deltaTime

		// Use last values if current is zero
		if rxBps == 0 {
			rxBps = lastNetIn
		}
		if txBps == 0 {
			txBps = lastNetOut
		}

		// Calculate percentage
		totalBps := rxBps + txBps
		usagePercent := (totalBps / 12500000) * 100

		// Update last values
		lastNetIn = rxBps
		lastNetOut = txBps
		lastPercentage = usagePercent

		// Save to file
		saveCurrentValues(rxBps, txBps, usagePercent)

		fmt.Printf("%s\t%.2f\t%.2f\t%.2f%%\n",
			time.Now().Format("15:04:05"),
			rxBps,
			txBps,
			usagePercent)

		// Update previous values
		prevNetIn, prevNetOut = currentNetIn, currentNetOut
		prevTime = time.Now()
	}
}
