package main

import (
	"fmt"
	"ralb_go_haproxy/utils"
	"time"
)

func main() {
	cfg := utils.LoadRalbEnv()
	utils.InitHTTPClient()

	var prevResults []utils.VMMetric
	var changeCount int = 0

	for {
		vms, err := utils.FetchVMs(cfg)
		if err != nil {
			fmt.Println("Error fetching VMs:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// var vmResults []utils.VMMetric
		var vmMetrics []utils.VMMetric
		for _, vm := range vms {
			if cfg.VMNames[vm.Name] && vm.Status == "running" {
				// vmResults = append(vmResults, utils.VMMetric{
				// 	Name:  vm.Name,
				// 	Score: utils.ResourceUsage(vm),
				// })
				vmMetrics = append(vmMetrics, utils.VMMetric{
					Name:      vm.Name,
					CPU:       vm.CPU,
					Memory:    (vm.Mem / vm.MaxMem),
					Bandwidth: (((vm.NetIn + vm.NetOut) / (1024 * 1024)) / 1000),
					Score:     utils.ResourceUsage(vm),
				})
			}
		}

		// Urutkan secara ascending dan beri/tentukan prioritas
		vmMetrics = utils.AscendingScoreSort(vmMetrics)

		// Bagian ini mengecek apakah hasil sebelumnya masih sama atau tidak
		// Bila sama maka jangan, bila tidak sama maka atau ada perubahan maka update
		same := len(vmMetrics) == len(prevResults)

		// utils.logVMsToCSV("logs", vmMetrics)

		if same {
			for i := range vmMetrics {
				if vmMetrics[i].Name != prevResults[i].Name || vmMetrics[i].Priority != prevResults[i].Priority {
					same = false
					break
				}
			}
		}

		if !same {
			changeCount++
			fmt.Printf("\n+---------------   DAFTAR VM   ----------------+\n")
			fmt.Printf("+--- NAMA VM ---+--- SKOR ---+--- PRIORITAS ---+\n")
			for _, usage := range vmMetrics {
				fmt.Printf("|    %s     |    %.2f    |      %d          |\n", usage.Name, usage.Score, usage.Priority)
			}
			fmt.Printf("+----------------------------------------------+\n")
			fmt.Printf("| Pembaruan ke-%d      |  %v   |\n", changeCount, time.Now().Format("2006/01/02 15:04:05"))
			fmt.Printf("+----------------------------------------------+\n")
			fmt.Printf("| RALB UPDATER :        | LOGGER : ")

			prevResults = make([]utils.VMMetric, len(vmMetrics))
			copy(prevResults, vmMetrics)

			if cfg.RalbUpdater == 1 {
				utils.ModifyHAProxy(cfg, vmMetrics)
				fmt.Println("UPDATED")
			} else {
				fmt.Println("UNUPDATED")
			}
		}

		if cfg.Logger == 1 {
			utils.CSVLogger("data", vmMetrics, prevResults, !same)
		}

		time.Sleep(time.Duration(cfg.FetchDelay) * time.Millisecond)
	}
}
