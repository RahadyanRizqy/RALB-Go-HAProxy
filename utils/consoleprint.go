package utils

import "fmt"

func ConsolePrint(stats map[string]VMStats, ranked map[string]VMPriority, netIfaceRate float64) {
	fmt.Println("\n== Detailed VM Stats ==")
	for name, stat := range stats {
		priority := ranked[name].Priority
		bwUsagePercent := (stat.BwUsage / netIfaceRate)
		memUsagePercent := stat.MemUsage
		cpuUsagePercent := stat.VM.CPU

		fmt.Printf("%s\tRx: %.2f\tTx: %.2f\tBW: %.2f%%\tCPU: %.2f%%\tMem: %.2f%%\tScore: %.6f\tPriority: %d\n",
			name,
			stat.Rates.Rx,
			stat.Rates.Tx,
			bwUsagePercent,
			cpuUsagePercent,
			memUsagePercent,
			cpuUsagePercent+memUsagePercent+bwUsagePercent,
			priority,
		)
	}
}
