package main

import (
	"fmt"
	"ralb_go_haproxy/utils"
	"time"
)

func main() {
	cfg := utils.LoadRalbEnv()
	utils.InitHTTPClient()

	var prevResults []utils.VMResult
	var changeCount int = 0

	for {
		vms, err := utils.FetchVMs(cfg)
		if err != nil {
			fmt.Println("Error fetching VMs:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var vmResults []utils.VMResult
		for _, vm := range vms {
			if cfg.VMNames[vm.Name] && vm.Status == "running" {
				vmResults = append(vmResults, utils.VMResult{
					Name:  vm.Name,
					Score: utils.ResourceUsage(vm),
				})
			}
		}

		// Urutkan secara ascending dan beri/tentukan prioritas
		vmResults = utils.AscendingSortPriority(vmResults)

		// Bagian ini mengecek apakah hasil sebelumnya masih sama atau tidak
		// Bila sama maka jangan, bila tidak sama maka atau ada perubahan maka update
		same := len(vmResults) == len(prevResults)
		if same {
			for i := range vmResults {
				if vmResults[i].Name != prevResults[i].Name || vmResults[i].Priority != prevResults[i].Priority {
					same = false
					break
				}
			}
		}

		if !same {
			changeCount++
			fmt.Printf("\n+---------------   DAFTAR VM   --------------+\n")
			fmt.Printf("+--- NAMA VM ---+---  SKOR ---+--- PRIORITY ---+\n")
			for _, usage := range vmResults {
				fmt.Printf("|    %s     |    %.2f     |      %d       |\n", usage.Name, usage.Score, usage.Priority)
			}
			fmt.Printf("+--------------------------------------------+\n")
			fmt.Printf("Perubahan ke-%d\n", changeCount)

			prevResults = make([]utils.VMResult, len(vmResults))
			copy(prevResults, vmResults)

			utils.ModifyHAProxy(cfg, vmResults)
		}

		time.Sleep(time.Duration(cfg.FetchDelay) * time.Millisecond)
	}
}
