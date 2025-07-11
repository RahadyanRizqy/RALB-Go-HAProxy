package funcs

import "ralb_go_haproxy/utils"

func CalcPreviousStats(vm utils.VM, delta float64, netIfaceRate float64, lastValidRates map[string]utils.ActiveRates, prevStats map[string]utils.VM, currentRates map[string]utils.ActiveRates) utils.VMStats {
	stats := utils.VMStats{VM: vm}

	// Calc net rates
	rxRate := lastValidRates[vm.Name].Rx
	txRate := lastValidRates[vm.Name].Tx

	if prev, ok := prevStats[vm.Name]; ok {
		// Only update rates if we have new non-zero values
		if vm.CumNetIn > prev.CumNetIn {
			newRx := float64(vm.CumNetIn-prev.CumNetIn) / delta
			if newRx > 0 {
				rxRate = newRx
			}
		}
		if vm.CumNetOut > prev.CumNetOut {
			newTx := float64(vm.CumNetOut-prev.CumNetOut) / delta
			if newTx > 0 {
				txRate = newTx
			}
		}
	}

	// Store current rates
	stats.Rates = utils.ActiveRates{Rx: rxRate, Tx: txRate}
	currentRates[vm.Name] = stats.Rates

	// Update last valid rates if we have non-zero values
	if rxRate > 0 || txRate > 0 {
		lastValidRates[vm.Name] = stats.Rates
	}

	// Calculate metrics
	stats.MemUsage = vm.Mem / float64(vm.MaxMem)
	stats.BwUsage = float64(rxRate+txRate) / netIfaceRate
	stats.Score = vm.CPU + stats.MemUsage + stats.BwUsage

	return stats
}

func UpdatePreviousState(prevStats map[string]utils.VM, stats map[string]utils.VMStats) {
	for name, stat := range stats {
		prevStats[name] = stat.VM
		// prevScores[name] = stat.Score
	}
}
