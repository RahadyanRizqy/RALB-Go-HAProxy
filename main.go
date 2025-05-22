package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"ralb_go_haproxy/funcs"
	"ralb_go_haproxy/utils"
	"time"
)

var client *http.Client

func InitHTTPClient() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func HasChanged(prev, current []utils.VMMetric) bool {
	if len(prev) != len(current) {
		return true
	}
	for i := range current {
		if current[i].Name != prev[i].Name || current[i].Priority != prev[i].Priority {
			return true
		}
	}
	return false
}

func main() {
	cfg := utils.LoadRalbEnv()
	InitHTTPClient()

	var prevResults []utils.VMMetric
	var updateCount int
	var logLine int = 1
	var scrapeCount int = 1

	for {
		rawVMs, err := funcs.FetchVMs(cfg, client)
		if err != nil {
			fmt.Printf("[%s] Error fetching VMs: %v\n", time.Now().Format("15:04:05"), err)
			time.Sleep(2 * time.Second)
			continue
		}

		vmMetrics := funcs.ExtractMetrics(cfg, rawVMs)
		changed := HasChanged(prevResults, vmMetrics)

		if changed {
			updateCount++
			funcs.ConsolePrint(vmMetrics, updateCount, cfg)

			prevResults = make([]utils.VMMetric, len(vmMetrics))
			copy(prevResults, vmMetrics)

			if cfg.RalbUpdater {
				funcs.UpdateHAProxy(cfg, vmMetrics)
			}
		}

		if cfg.Logger {
			utils.CSVLogger("data", vmMetrics, prevResults, updateCount, &logLine, scrapeCount)
		}

		scrapeCount++
		time.Sleep(time.Duration(cfg.FetchDelay) * time.Millisecond)
	}
}
