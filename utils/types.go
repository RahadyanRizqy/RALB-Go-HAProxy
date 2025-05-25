package utils

import "time"

type VM struct {
	Id        int     `json:"vmid"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	MaxMem    int     `json:"maxmem"`
	MaxCPU    int     `json:"maxcpu"`
	Mem       float64 `json:"mem"`
	CPU       float64 `json:"cpu"`
	CumNetIn  int     `json:"netin"`
	CumNetOut int     `json:"netout"`
}

type Response struct {
	Data []VM `json:"data"`
}

type VMPriority struct {
	Value    float64
	Priority int
}

type ActiveRates struct {
	Rx float64 // Receive rate in bytes/sec
	Tx float64 // Transmit rate in bytes/sec
}

type KV struct {
	Key   string
	Value float64
}

type VMStats struct {
	VM       VM
	Score    float64
	Rates    ActiveRates
	BwUsage  float64
	MemUsage float64
}

type VMLogs struct {
	VM        VM
	Score     float64
	RxRate    float64
	TxRate    float64
	MemUsage  float64
	BwUsage   float64
	Timestamp time.Time
}

type RalbEnv struct {
	APIToken             string
	PveAPIURL            string
	VMNames              map[string]bool
	HAProxyPath          string
	RalbUpdater          bool
	Logger               bool
	FetchDelay           int
	NetIfaceRate         float64
	ServerStart          bool
	ServerSuccessMessage string
	ServerErrorMessage   string
	ServerPort           int
}
