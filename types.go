package main

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

type VMResult struct {
	Name   string
	Score  float64
	Weight int
}

type RalbEnv struct {
	APIToken    string
	PveAPIURL   string
	VMNames     map[string]bool
	HAProxyPath string
	RalbStatus  int
	FetchDelay  int
}
