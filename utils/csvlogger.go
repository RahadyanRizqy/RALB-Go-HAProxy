package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func InitCSV() string {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return ""
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("data/vm_metrics_%d.csv", timestamp)

	// Create and open the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
	}
	defer file.Close()

	// Write headers
	writer := csv.NewWriter(file)
	headers := []string{
		"no", "fetch", "update", "vm_id", "vm_name", "cpu_usage", "max_cpu",
		"mem_usage", "max_mem", "cum_netin", "cum_netout", "rate_netin",
		"rate_netout", "bw_usage", "max_bw", "score", "priority",
		"unix_timestamp", "timestamp",
	}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
		return ""
	}
	writer.Flush()

	return filename
}

func StoreCSV(
	filename string,
	logLine *int,
	fetchCount int,
	updateCount int,
	unix_timestamp int64,
	timestamp string,
	stats map[string]VMStats,
	ranked map[string]VMPriority,
	netIfaceRate float64,
) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for name, stat := range stats {
		priority := ranked[name].Priority
		bw_usage := (stat.BwUsage / netIfaceRate)
		mem_usage := (stat.MemUsage)
		cpu_usage := stat.VM.CPU

		record := []string{
			strconv.Itoa(*logLine),
			strconv.Itoa(fetchCount),
			strconv.Itoa(updateCount),
			strconv.Itoa(stat.VM.Id),
			stat.VM.Name,
			strconv.FormatFloat(cpu_usage, 'f', -1, 64),
			strconv.Itoa(stat.VM.MaxCPU),
			strconv.FormatFloat(mem_usage, 'f', -1, 64),
			strconv.Itoa(stat.VM.MaxMem),
			strconv.Itoa(stat.VM.CumNetIn),
			strconv.Itoa(stat.VM.CumNetOut),
			strconv.FormatFloat(stat.Rates.Rx, 'f', -1, 64),
			strconv.FormatFloat(stat.Rates.Tx, 'f', -1, 64),
			strconv.FormatFloat(bw_usage, 'f', -1, 64),
			strconv.FormatFloat(netIfaceRate, 'f', -1, 64),
			strconv.FormatFloat(cpu_usage+mem_usage+bw_usage, 'f', -1, 64),
			strconv.Itoa(priority),
			strconv.FormatInt(unix_timestamp, 10),
			timestamp,
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
		*logLine++
	}
	return nil
}
