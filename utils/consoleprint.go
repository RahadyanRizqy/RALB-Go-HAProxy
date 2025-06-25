package utils

import (
	"fmt"
)

func ConsolePrint(cfg RalbEnv, stats map[string]VMStats, ranked map[string]VMPriority, netIfaceRate float64) error {
	if cfg.ConsolePrint {
		fmt.Println("\n== Detailed VM Stats ==")
		for name, stat := range stats {
			fmt.Printf("%s\tRx: %.2f\tTx: %.2f\tBW: %.2f%%\tCPU: %.2f%%\tMem: %.2f%%\tScore: %.6f\tPriority: %d\n",
				name,
				stat.Rates.Rx,
				stat.Rates.Tx,
				stat.BwUsage,
				stat.VM.CPU,
				stat.MemUsage,
				stat.Score,
				ranked[name].Priority,
			)
		}
	}

	return nil
}
