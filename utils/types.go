package utils

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
	Weight   int
}

type VMRank struct {
	Value    float64
	Priority int
	Weight   int
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

type RalbEnv struct {
	APIToken             string
	PveAPIURL            string
	HAProxySock          string
	HAProxyBackend       string
	HAProxyWeight        int
	VMNames              map[string]bool
	VMIPs                map[string]bool
	VMMap                map[string]string
	HAProxyPath          string
	RalbUpdater          bool
	UpdateNotify         bool
	HAProxySetWeight     bool
	ConsolePrint         bool
	Logger               bool
	FetchDelay           int
	NetIfaceRate         float64
	ServerStart          bool
	ServerSuccessMessage string
	ServerErrorMessage   string
	ServerPort           int
	Strict               bool
}
