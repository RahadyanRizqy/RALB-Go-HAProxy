package funcs

import (
	"math"
	"ralb_go_haproxy/utils"
	"sort"
)

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

func CalcScorePriorityWeight(stats map[string]utils.VMStats, cfg utils.RalbEnv) map[string]utils.VMRank {
	var sorted []utils.KV
	for name, stat := range stats {
		sorted = append(sorted, utils.KV{Key: name, Value: stat.Score})
	}

	// Sort by value
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

	n := len(result)

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
	_result := make(map[string]utils.VMRank)
	for name, vm := range result {
		_result[name] = utils.VMRank{
			Value:    vm.Value,
			Priority: vm.Priority,
			Weight:   priorityToWeight[vm.Priority],
		}
	}

	return _result
}
