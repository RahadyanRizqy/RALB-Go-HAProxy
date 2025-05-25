package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"ralb_go_haproxy/funcs"
	"ralb_go_haproxy/utils"
	"sort"
	"time"
)

var (
	prevStats      = make(map[string]utils.VM)
	prevScores     = make(map[string]float64)
	activeRates    = make(map[string]utils.ActiveRates)
	lastValidRates = make(map[string]utils.ActiveRates)
	client         *http.Client
	fetchCount     int
	updateCount    int
	logLine        int = 1
)

func InitClient() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func Start() {
	InitClient()
	cfg := utils.LoadRalbEnv()
	csvFileName := utils.InitCSV()
	prevTime := time.Now()

	for {
		time.Sleep(1 * time.Second)
		now := time.Now()
		delta := now.Sub(prevTime).Seconds()
		fetchCount++

		stats, err := funcs.FetchVMs(cfg, client)
		if err != nil {
			fmt.Printf("Polling error: %v\n", err)
			continue
		}

		// Process VM stats and calculate metrics
		currentStats := make(map[string]utils.VMStats)
		scoreChanged := false

		for _, vm := range stats {
			if !cfg.VMNames[vm.Name] {
				continue
			}

			stats := processVMStats(vm, delta)
			currentStats[vm.Name] = stats

			// Check for score changes
			if prevScore, exists := prevScores[vm.Name]; exists {
				if stats.Score != prevScore {
					scoreChanged = true
				}
			} else {
				scoreChanged = true
			}
		}

		// Print notification if scores changed

		if scoreChanged {
			fmt.Printf("\n=== NEW SCORE DETECTED (Update #%d) ===\n", updateCount)
			updateCount++
		}
		// Print current metrics
		fmt.Printf("\n[%s] Fetch #%d\n", now.Format("15:04:05"), fetchCount)

		// Sort and rank VMs by score
		rankedVMs := rankVMsByScore(currentStats)

		// Print full VM stats
		// logLine++
		utils.ConsolePrint(currentStats, rankedVMs, cfg.NetIfaceRate)
		utils.StoreCSV(
			csvFileName,
			&logLine,
			fetchCount,
			updateCount,
			now.Unix(),
			now.Format("2006-01-02 15:04:05"),
			currentStats,
			rankedVMs,
			cfg.NetIfaceRate)

		// Update previous state
		updatePreviousState(currentStats)
		prevTime = now
	}
}

func processVMStats(vm utils.VM, delta float64) utils.VMStats {
	stats := utils.VMStats{VM: vm}

	// Calculate network rates
	rxRate := lastValidRates[vm.Name].Rx // Start with last valid rate
	txRate := lastValidRates[vm.Name].Tx

	if prev, ok := prevStats[vm.Name]; ok {
		// Only update rates if we have new non-zero values
		if vm.CumNetIn > prev.CumNetIn {
			newRx := float64(vm.CumNetIn-prev.CumNetIn) / delta
			if newRx > 0 {
				rxRate = newRx
			}
		}
		if vm.CumNetOut > prev.CumNetOut {
			newTx := float64(vm.CumNetOut-prev.CumNetOut) / delta
			if newTx > 0 {
				txRate = newTx
			}
		}
	}

	// Store current rates
	stats.Rates = utils.ActiveRates{Rx: rxRate, Tx: txRate}
	activeRates[vm.Name] = stats.Rates

	// Update last valid rates if we have non-zero values
	if rxRate > 0 || txRate > 0 {
		lastValidRates[vm.Name] = stats.Rates
	}

	// Calculate metrics
	stats.MemUsage = vm.Mem / float64(vm.MaxMem)
	stats.BwUsage = (rxRate + txRate)
	stats.Score = vm.CPU + stats.MemUsage + stats.BwUsage

	return stats
}

func rankVMsByScore(stats map[string]utils.VMStats) map[string]utils.VMPriority {
	// Convert to slice for sorting
	var sorted []utils.KV
	for name, stat := range stats {
		sorted = append(sorted, utils.KV{Key: name, Value: stat.Score})
	}

	// Sort by value
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	// Build ranked map
	result := make(map[string]utils.VMPriority)
	for i, item := range sorted {
		result[item.Key] = utils.VMPriority{
			Value:    item.Value,
			Priority: i + 1,
		}
	}

	return result
}

func updatePreviousState(stats map[string]utils.VMStats) {
	for name, stat := range stats {
		prevStats[name] = stat.VM
		prevScores[name] = stat.Score
	}
}
