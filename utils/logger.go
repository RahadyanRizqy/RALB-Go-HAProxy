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
func CSVLogger(dir string, metrics []VMMetric, prev []VMMetric, update int, logLine *int, scrapeCount int) {
	// Ensure logs directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating log directory:", err)
		return
	}

	filename := filepath.Join(dir, "vm_metrics.csv")
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
		writer.Write([]string{"no", "scrape", "update", "vm_name", "cpu_usage", "mem_usage", "bw_usage", "score", "priority", "timestamp"})
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Write current metrics
	for _, vm := range metrics {
		record := []string{
			strconv.Itoa(*logLine),
			strconv.Itoa(scrapeCount),
			fmt.Sprintf("%d", update),
			vm.Name,
			fmt.Sprintf("%f", vm.CPU),
			fmt.Sprintf("%f", vm.Memory),
			fmt.Sprintf("%f", vm.Bandwidth),
			fmt.Sprintf("%f", vm.Score),
			fmt.Sprintf("%d", vm.Priority),
			timestamp,
		}
		writer.Write(record)
		*logLine++
	}
}
