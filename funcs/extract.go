package funcs

import "ralb_go_haproxy/utils"

func ExtractMetrics(cfg utils.RalbEnv, vms []utils.VM) []utils.VMMetric {
	var result []utils.VMMetric
	for _, vm := range vms {
		if cfg.VMNames[vm.Name] && vm.Status == "running" {
			result = append(result, utils.VMMetric{
				VM:             vm,
				BandwidthUsage: 0,
				BandwidthRate:  0,
				Score:          ResourceUsage(vm),
			})
		}
	}
	return AscendingScoreSort(result)
}
