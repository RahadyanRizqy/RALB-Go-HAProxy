package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func ModifyHAProxy(cfg RalbEnv, vmResults []VMResult) {
	input, err := os.ReadFile(cfg.HAProxyPath)
	if err != nil {
		log.Printf("Failed to read HAProxy config: %v\n", err)
		return
	}
	lines := strings.Split(string(input), "\n")
	weightMap := make(map[string]int)

	for _, vm := range vmResults {
		weightMap[vm.Name] = vm.Weight
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "server") {
			fields := strings.Fields(trimmed)
			if len(fields) >= 8 {
				vmName := fields[1]
				if weight, ok := weightMap[vmName]; ok {
					for j, f := range fields {
						if f == "weight" && j+1 < len(fields) {
							fields[j+1] = strconv.Itoa(weight)
							leading := line[:strings.Index(line, "server")]
							lines[i] = leading + strings.Join(fields, " ")
							break
						}
					}
				}
			}
		}
	}

	output := strings.Join(lines, "\n")
	if err := os.WriteFile(cfg.HAProxyPath, []byte(output), 0644); err != nil {
		log.Printf("Failed to write updated HAProxy config: %v\n", err)
	} else {
		log.Println("HAProxy config successfully updated.")
	}
}
