package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	cfg := LoadRalbEnv()
	InitHTTPClient()

	var lastWeights []VMResult
	var changeCount int

	for {
		vms, err := FetchVMs(cfg)
		if err != nil {
			fmt.Println("Error fetching VMs:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var vmResults []VMResult
		for _, vm := range vms {
			if cfg.VMNames[vm.Name] && vm.Status == "running" {
				vmResults = append(vmResults, VMResult{
					Name:  vm.Name,
					Score: ResourceUsage(vm),
				})
			}
		}

		// sort and assign weights
		sort.Slice(vmResults, func(i, j int) bool {
			return vmResults[i].Score < vmResults[j].Score
		})
		for i := range vmResults {
			vmResults[i].Weight = len(vmResults) - i
		}

		// Compare weights to last run
		same := len(vmResults) == len(lastWeights)
		if same {
			for i := range vmResults {
				if vmResults[i].Name != lastWeights[i].Name || vmResults[i].Weight != lastWeights[i].Weight {
					same = false
					break
				}
			}
		}

		// Print and apply if there's a change
		if !same {
			changeCount++
			fmt.Printf("\n+---------------   DAFTAR VM   --------------+\n")
			fmt.Printf("+--- NAMA VM ---+--- SKOR  ---+--- WEIGHT ---+\n")
			for _, usage := range vmResults {
				fmt.Printf("|    %s     |    %.2f     |      %d       |\n", usage.Name, usage.Score, usage.Weight)
			}
			fmt.Printf("+--------------------------------------------+\n")
			fmt.Printf("Change (th): %d\n", changeCount)

			lastWeights = make([]VMResult, len(vmResults))
			copy(lastWeights, vmResults)

			ModifyHAProxy(cfg, vmResults)
		}

		time.Sleep(time.Duration(cfg.FetchDelay) * time.Millisecond)
	}
}
