package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Updated function to accept previous metrics and a flag whether it changed
func CSVLogger(logLine *int, updateCount int, scrapeCount int, metrics []VMMetric) {
	// Ensure logs directory exists
	const dir string = "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}

	filename := filepath.Join(dir, "vm_results.csv")
	isNewFile := false

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		isNewFile = true
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if isNewFile {
		writer.Write([]string{"no", "scrape", "update", "vm_name", "vm_id", "max_cpu", "cpu_usage", "max_mem", "mem_usage", "netin", "netout", "bw_rate", "bw_usage", "score", "priority", "timestamp"})
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Write current metrics
	for _, vm := range metrics {
		record := []string{
			strconv.Itoa(*logLine),
			strconv.Itoa(scrapeCount),
			fmt.Sprintf("%d", updateCount),

			vm.Name,
			strconv.Itoa(vm.Id),
			fmt.Sprintf("%f", vm.MaxCPU),
			fmt.Sprintf("%f", vm.CPU),
			fmt.Sprintf("%f", vm.MaxMem),
			fmt.Sprintf("%f", vm.Mem),
			fmt.Sprintf("%f", vm.NetIn),
			fmt.Sprintf("%f", vm.NetOut),
			fmt.Sprintf("%f", vm.BandwidthRate),
			fmt.Sprintf("%f", vm.BandwidthUsage),
			fmt.Sprintf("%f", vm.Score),
			fmt.Sprintf("%d", vm.Priority),

			timestamp,
		}
		writer.Write(record)
		*logLine++
	}
}
