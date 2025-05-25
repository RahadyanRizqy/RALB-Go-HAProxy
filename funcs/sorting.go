package funcs

// import (
// 	"math"
// 	"ralb_go_haproxy/utils"
// 	"sort"
// )

// func ResourceUsage(vm utils.VM) float64 {
// 	usage := vm.CPU + (vm.Mem / vm.MaxMem) + (((vm.NetIn + vm.NetOut) / (1024 * 1024)) / 1000)
// 	return math.Round(usage*100) / 100
// }

// func AscendingScoreSort(results []utils.VMMetric) []utils.VMMetric {
// 	sort.Slice(results, func(i, j int) bool {
// 		return results[i].Score < results[j].Score
// 	})

// 	// Assign priority (1 = highest priority)
// 	for i := range results {
// 		results[i].Priority = len(results) - (len(results) - i) + 1
// 	}

// 	return results
// }
