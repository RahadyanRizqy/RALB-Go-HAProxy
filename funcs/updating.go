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

func AssignWeightByPriority(ranked map[string]utils.VMPriority, cfg utils.RalbEnv) map[string]utils.VMPriority {
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

// func UpdateHAProxy(cfg utils.RalbEnv, ranked map[string]utils.VMPriority) error {
// 	// Combine VM names and IPs
// 	combinedVMs := CombineVMs(cfg.VMNames, cfg.VMIPs)
// 	fmt.Println("Combined VM map:", combinedVMs)

// 	if cfg.RalbUpdater {
// 		if len(cfg.VMNames) != len(cfg.VMIPs) {
// 			return fmt.Errorf("VM_NAMES and VM_IPS must be the same length")
// 		}

// 		// Build ordered list by priority
// 		type kv struct {
// 			Name     string
// 			Priority int
// 		}
// 		rankSlice := make([]kv, 0, len(ranked))
// 		for name, v := range ranked {
// 			rankSlice = append(rankSlice, kv{name, v.Priority})
// 		}
// 		sort.Slice(rankSlice, func(i, j int) bool {
// 			return rankSlice[i].Priority < rankSlice[j].Priority
// 		})

// 		// Create reordered list based on priority
// 		type vmPair struct {
// 			Name string
// 			IP   string
// 		}
// 		reordered := make([]vmPair, 0, len(rankSlice))
// 		for _, item := range rankSlice {
// 			ip, ok := combinedVMs[item.Name]
// 			if !ok {
// 				fmt.Printf("Warning: VM name %s not found in combinedVMs\n", item.Name)
// 				continue
// 			}
// 			reordered = append(reordered, vmPair{Name: item.Name, IP: ip})
// 		}

// 		// Read HAProxy config
// 		contentBytes, err := os.ReadFile(cfg.HAProxyPath)
// 		if err != nil {
// 			return fmt.Errorf("failed to read haproxy config: %v", err)
// 		}
// 		lines := strings.Split(string(contentBytes), "\n")

// 		updatedLines := make([]string, 0, len(lines))
// 		inBackend := false
// 		inserted := false

// 		for i := 0; i < len(lines); i++ {
// 			line := lines[i]
// 			trimmed := strings.TrimSpace(line)

// 			if strings.HasPrefix(trimmed, "backend web_servers") {
// 				inBackend = true
// 				updatedLines = append(updatedLines, line)
// 				continue
// 			}

// 			if inBackend {
// 				if strings.HasPrefix(trimmed, "option httpchk") && !inserted {
// 					updatedLines = append(updatedLines, line)

// 					for idx, p := range reordered {
// 						serverLine := fmt.Sprintf("        server %s %s:80 check inter 2s rise 3 fall 2", p.Name, p.IP)
// 						if idx > 0 {
// 							serverLine += " backup"
// 						}
// 						updatedLines = append(updatedLines, serverLine)
// 					}

// 					// Skip old server lines
// 					for i+1 < len(lines) {
// 						next := strings.TrimSpace(lines[i+1])
// 						if strings.HasPrefix(next, "server ") {
// 							i++
// 						} else {
// 							break
// 						}
// 					}

// 					inserted = true
// 					continue
// 				}

// 				if trimmed == "" || strings.HasPrefix(trimmed, "frontend ") || strings.HasPrefix(trimmed, "backend ") {
// 					inBackend = false
// 				}
// 			}

// 			updatedLines = append(updatedLines, line)
// 		}

// 		if !inserted {
// 			return fmt.Errorf("failed to update HAProxy weightg: did not find 'option httpchk' to insert servers")
// 		}

// 		err = os.WriteFile(cfg.HAProxyPath, []byte(strings.Join(updatedLines, "\n")), 0644)
// 		if err != nil {
// 			return fmt.Errorf("failed to write updated config: %v", err)
// 		}
// 	}

// 	// Reload HAProxy if enabled
// 	if cfg.HAProxySetWeight {
// 		cmd := exec.Command("bash", "-c", fmt.Sprintf("haproxy -f %s -sf $(cat /var/run/haproxy.pid)", cfg.HAProxyPath))
// 		output, err := cmd.CombinedOutput()
// 		if err != nil {
// 			return fmt.Errorf("failed to reload haproxy: %v, output: %s", err, string(output))
// 		}
// 		fmt.Println("HAProxy reloaded successfully.")
// 	}

// 	return nil
// }

func UpdateHAProxy(cfg utils.RalbEnv, ranked map[string]utils.VMPriority) error {
	if !cfg.HAProxySetWeight {
		return nil
	}

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

			fmt.Printf("Set weight untuk %s: %d\n", vmName, weight)
		}
	}

	return nil
}
