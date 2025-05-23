package funcs

import (
	"fmt"
	"ralb_go_haproxy/utils"
	"time"
)

func ConsolePrint(vmMetrics []utils.VMMetric, updateCount int, cfg utils.RalbEnv) {
	fmt.Printf("\n+---------------   DAFTAR VM   ----------------+\n")
	fmt.Printf("+--- NAMA VM ---+--- SKOR ---+--- PRIORITAS ---+\n")
	for _, usage := range vmMetrics {
		fmt.Printf("| %-13s | %-9.2f | %-16d |\n", usage.Name, usage.Score, usage.Priority)
	}
	fmt.Printf("+----------------------------------------------+\n")
	fmt.Printf("| %-16s: %-27d|\n", "PEMBARUAN KE", updateCount)
	fmt.Printf("| %-16s: %-27s|\n", "TIMESTAMP", time.Now().Format("2006/01/02 15:04:05"))
	fmt.Printf("| %-16s: %-27t|\n", "RALB UPDATER", cfg.RalbUpdater)
	fmt.Printf("| %-16s: %-27t|\n", "LOGGER", cfg.Logger)
	fmt.Printf("| %-16s: %-27t|\n", "SERVER START", cfg.ServerStart)
	fmt.Printf("| %-16s: %-27t|\n", "SERVER PORT", cfg.ServerStart)
	fmt.Printf("+----------------------------------------------+\n")
}
