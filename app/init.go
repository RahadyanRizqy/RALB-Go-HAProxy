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
	fmt.Println("RALB Started!")
	InitClient()
	cfg := utils.LoadRalbEnv()
	csvFileName := utils.InitCSV(cfg)
	prevTime := time.Now()

	for {
		time.Sleep(time.Duration(cfg.FetchDelay) * time.Millisecond)
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

			stats := funcs.PreviousStats(vm, delta, cfg.NetIfaceRate, lastValidRates, prevStats, activeRates)
			currentStats[vm.Name] = stats

			// Check for score changes
			if prevScore, exists := prevScores[vm.Name]; exists {
				if stats.Score != prevScore {
					scoreChanged = true

				}
			}
		}

		rankedVMs := funcs.ScorePriority(currentStats)
		rankedWithWeight := funcs.AssignWeightByPriority(rankedVMs, cfg)
		// Print notification if scores changed
		if scoreChanged {
			updateCount++
			if cfg.UpdateNotify {
				fmt.Printf("\n=== NEW SCORE DETECTED (Update #%d) ===\n", updateCount)
			}
			funcs.UpdateHAProxy(cfg, rankedWithWeight)
		}

		// // Print current metrics
		// fmt.Printf("\n[%s] Fetch #%d\n", now.Format("15:04:05"), fetchCount)

		// Sort and rank VMs by score

		// Print full VM stats
		utils.ConsolePrint(cfg, currentStats, rankedVMs, cfg.NetIfaceRate)
		utils.StoreCSV(
			cfg,
			csvFileName,
			&logLine,
			fetchCount,
			updateCount,
			now.Unix(),
			now.Format("2006-01-02 15:04:05"),
			currentStats,
			rankedWithWeight,
			cfg.NetIfaceRate)

		// Update previous state
		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)
		prevTime = now
	}
}
