package utils

import "time"

type VM struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Status string  `json:"status"`
	VMID   int     `json:"vmid"`
	MaxMem float64 `json:"maxmem"`
	MaxCPU float64 `json:"maxcpu"`
	Mem    float64 `json:"mem"`
	CPU    float64 `json:"cpu"`
	NetIn  float64 `json:"netin"`
	NetOut float64 `json:"netout"`
}

type Response struct {
	Data []VM `json:"data"`
}

type BandwidthTracker struct {
	lastGoodRX float64
	lastGoodTX float64
	prevNetIn  float64
	prevNetOut float64
	prevTime   time.Time
}

type VMMetric struct {
	Name      string
	CPU       float64
	Memory    float64
	Bandwidth float64
	Score     float64
	Priority  int
}

type RalbEnv struct {
	APIToken    string
	PveAPIURL   string
	VMNames     map[string]bool
	HAProxyPath string
	RalbUpdater bool
	Logger      bool
	RunServer   bool
	FetchDelay  int
}
