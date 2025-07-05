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
	prevWeights    = make(map[string]int)
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
		delta := now.Sub(prevTime).Seconds() // To differentiate bandwidth rate
		fetchCount++

		/*
			FetchStats() to fetch VM stats from Proxmox VE API for logging and RALB
		*/
		stats, err := funcs.FetchStats(cfg, client)
		if err != nil {
			fmt.Printf("Polling error: %v\n", err)
			continue
		}

		/*
			Calculate Previous Stats to calculate rate between fetches by calculating it's delta
		*/
		currentStats := make(map[string]utils.VMStats)
		for _, vm := range stats {
			if !cfg.VMNames[vm.Name] {
				continue
			}
			currentStats[vm.Name] = funcs.CalcPreviousStats(vm, delta, cfg.NetIfaceRate, lastValidRates, prevStats, activeRates)
		}

		/*
			Ranked Result
		*/
		rankedResult := funcs.CalcScorePriorityWeight(currentStats, cfg)

		/*
			At least 1 VM is updated
		*/
		var weightChanged bool = false

		for name, current := range rankedResult {
			if prev, ok := prevWeights[name]; !ok || prev != current.Weight {
				weightChanged = true
				break
			}
		}

		/*
			if new weight detected change the weight
		*/
		if weightChanged {
			fmt.Println(rankedResult)
			funcs.SetWeight(rankedResult, cfg)
			updateCount++
			fmt.Println("Update ke-", updateCount)
		}

		/*
			update the previous weight to differentiate the previous and current
		*/
		for name, current := range rankedResult {
			prevWeights[name] = current.Weight
		}

		/*
			Console Print
		*/
		utils.ConsolePrint(currentStats, rankedResult, cfg)

		/*
			Log the Result to CSV
		*/
		utils.StoreCSV(
			cfg,
			csvFileName,
			&logLine,
			fetchCount,
			updateCount,
			now.Unix(),
			now.Format("2006-01-02 15:04:05"),
			currentStats,
			rankedResult,
			cfg.NetIfaceRate)

		/*
			Update previous stats
		*/
		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)
		prevTime = now
	}
}
