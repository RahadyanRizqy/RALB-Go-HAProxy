package funcs

import (
	"fmt"
	"os/exec"
	"ralb_go_haproxy/utils"
)

func SetWeight(ranked map[string]utils.VMRank, cfg utils.RalbEnv) error {
	// Ambil backend dan path sock
	backend := cfg.HAProxyBackend
	sockPath := cfg.HAProxySock

	for vmName, data := range ranked {
		weight := data.Weight

		// Perintah shell untuk mengatur bobot
		cmdStr := fmt.Sprintf(`echo "set weight %s/%s %d" | socat stdio %s`, backend, vmName, weight, sockPath)

		if cfg.RalbUpdater {
			cmd := exec.Command("bash", "-c", cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Gagal set weight untuk %s: %v\nOutput: %s\n", vmName, err, string(output))
				continue
			}
			if cfg.ConsolePrint {
				fmt.Printf("Set weight untuk %s: %d\n", vmName, weight)
			}
		}
	}

	return nil
}
