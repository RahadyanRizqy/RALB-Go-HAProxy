package funcs

import (
	"ralb_go_haproxy/utils"
	"sort"
)

func ScorePriority(stats map[string]utils.VMStats) map[string]utils.VMPriority {
	var sorted []utils.KV
	for name, stat := range stats {
		sorted = append(sorted, utils.KV{Key: name, Value: stat.Score})
	}

	// Sort by value
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	// Build ranked map
	result := make(map[string]utils.VMPriority)
	for i, item := range sorted {
		result[item.Key] = utils.VMPriority{
			Value:    item.Value,
			Priority: i + 1,
		}
	}

	return result
}
