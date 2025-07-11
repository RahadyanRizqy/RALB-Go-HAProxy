package funcs

import (
	"math"
	"ralb_go_haproxy/utils"
	"sort"
)

func DistributeWeights(_result map[string]utils.VMRank, weightTotal int) map[int]int {
	n := len(_result)

	// Buat array bobot berdasarkan deret (1,2,...,n)
	base := make([]int, len(_result))
	for i := 0; i < n; i++ {
		base[i] = i + 1
	}

	sum := 0
	for _, v := range base {
		sum += v
	}

	result := make([]int, len(base))
	for i, val := range base {
		ratio := float64(val) / float64(sum)
		result[i] = int(math.Round(ratio * float64(weightTotal)))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(result)))

	// Mapping priority ke bobot
	distributedWeights := make(map[int]int)
	for i, w := range result {
		distributedWeights[i+1] = w // prioritas 1 → bobot terbesar
	}
	return distributedWeights
}

func CalcScorePriorityWeight(stats map[string]utils.VMStats, cfg utils.RalbEnv) map[string]utils.VMRank {
	var sorted []utils.KV
	for name, stat := range stats {
		sorted = append(sorted, utils.KV{Key: name, Value: stat.Score})
	}

	// Sort by Value
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	// Build ranked map
	result := make(map[string]utils.VMRank)
	for i, item := range sorted {
		result[item.Key] = utils.VMRank{
			Value:    item.Value,
			Priority: i + 1,
		}
	}

	// Hitung bobot proporsional dari deret ke totalWeight
	weights := DistributeWeights(result, cfg.HAProxyWeight)

	// Bangun hasil akhir
	_result := make(map[string]utils.VMRank)
	for name, vm := range result {
		_result[name] = utils.VMRank{
			Value:    vm.Value,
			Priority: vm.Priority,
			Weight:   weights[vm.Priority],
		}
	}

	return _result
}
