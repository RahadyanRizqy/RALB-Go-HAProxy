package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"ralb_go_haproxy/funcs"
	"ralb_go_haproxy/utils"
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

			stats := funcs.PreviousStats(vm, delta, lastValidRates, prevStats, activeRates)
			currentStats[vm.Name] = stats

			// Check for score changes
			if prevScore, exists := prevScores[vm.Name]; exists {
				if stats.Score != prevScore {
					scoreChanged = true

				}
			}
			// else {
			// 	// scoreChanged = true
			// 	updateCount++
			// }
		}

		// Print notification if scores changed
		if scoreChanged {
			updateCount++
			fmt.Printf("\n=== NEW SCORE DETECTED (Update #%d) ===\n", updateCount)
		}

		// Print current metrics
		fmt.Printf("\n[%s] Fetch #%d\n", now.Format("15:04:05"), fetchCount)

		// Sort and rank VMs by score
		rankedVMs := funcs.ScorePriority(currentStats)

		// Print full VM stats
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
		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)
		prevTime = now
	}
}
