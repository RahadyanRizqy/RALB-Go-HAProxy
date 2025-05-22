package funcs

import "ralb_go_haproxy/utils"

func ExtractMetrics(cfg utils.RalbEnv, vms []utils.VM) []utils.VMMetric {
	var result []utils.VMMetric
	for _, vm := range vms {
		if cfg.VMNames[vm.Name] && vm.Status == "running" {
			result = append(result, utils.VMMetric{
				Name:      vm.Name,
				CPU:       vm.CPU,
				Memory:    vm.Mem / vm.MaxMem,
				Bandwidth: (vm.NetIn + vm.NetOut) / (1024 * 1024) / 1000,
				Score:     ResourceUsage(vm),
			})
		}
	}
	return AscendingScoreSort(result)
}
