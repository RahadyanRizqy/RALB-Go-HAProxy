package funcs

import (
	"fmt"
	"math"
	"os/exec"
	"ralb_go_haproxy/utils"
	"sort"
)

func CombineVMs(vmNames, vmIPs map[string]bool) map[string]string {
	result := make(map[string]string)

	// Convert keys to slices
	names := make([]string, 0, len(vmNames))
	for name := range vmNames {
		names = append(names, name)
	}
	sort.Strings(names) // Sort alphabetically

	ips := make([]string, 0, len(vmIPs))
	for ip := range vmIPs {
		ips = append(ips, ip)
	}
	sort.Strings(ips) // Sort numerically/alphabetically

	// Pair by index
	for i := 0; i < len(names) && i < len(ips); i++ {
		result[names[i]] = ips[i]
	}

	return result
}

func Sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}

func DistributeWeights(arr []int, weightTotal int) []int {
	sum := Sum(arr)
	result := make([]int, len(arr))
	for i, val := range arr {
		ratio := float64(val) / float64(sum)
		result[i] = int(math.Round(ratio * float64(weightTotal)))
	}
	return result
}

func WeightAssignment(ranked map[string]utils.VMPriority, cfg utils.RalbEnv) map[string]utils.VMPriority {
	n := len(ranked)

	// Buat array bobot berdasarkan deret (1,2,...,n)
	base := make([]int, n)
	for i := 0; i < n; i++ {
		base[i] = i + 1
	}

	// Hitung bobot proporsional dari deret ke totalWeight
	weights := DistributeWeights(base, cfg.HAProxyWeight)
	sort.Sort(sort.Reverse(sort.IntSlice(weights))) // urut dari besar ke kecil

	// Mapping priority ke bobot
	priorityToWeight := make(map[int]int)
	for i, w := range weights {
		priorityToWeight[i+1] = w // prioritas 1 → bobot terbesar
	}

	// Bangun hasil akhir
	result := make(map[string]utils.VMPriority)
	for name, vm := range ranked {
		result[name] = utils.VMPriority{
			Value:    vm.Value,
			Priority: vm.Priority,
			Weight:   priorityToWeight[vm.Priority],
		}
	}

	return result
}

func ChangeWeight(cfg utils.RalbEnv, ranked map[string]utils.VMPriority) error {
	// Ambil backend dan path sock
	backend := cfg.HAProxyBackend
	sockPath := cfg.HAProxySock

	for vmName, data := range ranked {
		weight := data.Weight

		// Perintah shell untuk mengatur bobot
		cmdStr := fmt.Sprintf(`echo "set weight %s/%s %d" | socat stdio %s`, backend, vmName, weight, sockPath)

		if cfg.RalbUpdater {
			cmd := exec.Command("bash", "-c", cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Gagal set weight untuk %s: %v\nOutput: %s\n", vmName, err, string(output))
				continue
			}
			if cfg.ConsolePrint {
				fmt.Printf("Set weight untuk %s: %d\n", vmName, weight)
			}
		}
	}

	return nil
}
