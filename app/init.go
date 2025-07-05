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

	iter := 1
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
		currentRes := funcs.CalcScorePriorityWeight(currentStats, cfg)

		/*
			Strict or Loose
		*/
		if cfg.Strict { // Strict means new weight of each VM must different from previous one
			validate1 := funcs.AllWeightValidation(currentRes, prevWeights)

			if validate1 {
				updateCount++
				fmt.Printf("✅ UPDATE COUNT %d ITER COUNT %d\n", updateCount, iter)
				funcs.SetWeight(currentRes, cfg)
				utils.ConsolePrint(currentStats, currentRes, cfg)
				for name, info := range currentRes {
					prevWeights[name] = info.Weight // update previous
				}
			}
		} else { // Loose means new weight of each VM has swapped not all but some
			validate2 := funcs.SomeWeightValidation(currentRes, prevWeights)

			if validate2 {
				updateCount++
				fmt.Printf("✅ UPDATE COUNT %d ITER COUNT %d\n", updateCount, iter)
				funcs.SetWeight(currentRes, cfg)
				utils.ConsolePrint(currentStats, currentRes, cfg)
				for name, info := range currentRes {
					prevWeights[name] = info.Weight // update previous
				}
			}
		}

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
			currentRes,
			cfg.NetIfaceRate)

		/*
			Update previous stats
		*/
		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)
		prevTime = now

		iter++
	}
}
