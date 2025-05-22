package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func parseVMMap(env string) map[string]bool {
	result := make(map[string]bool)
	for _, vm := range strings.Split(env, ",") {
		if trimmed := strings.TrimSpace(vm); trimmed != "" {
			result[trimmed] = true
		}
	}
	return result
}

func LoadRalbEnv() RalbEnv {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	ralb, _ := strconv.Atoi(os.Getenv("RALB_STATUS"))
	delay, err := strconv.Atoi(os.Getenv("FETCH_DELAY"))
	if err != nil {
		delay = 1000
	}

	logger, _ := strconv.Atoi(os.Getenv("LOGGER"))

	return RalbEnv{
		APIToken:    os.Getenv("API_TOKEN"),
		PveAPIURL:   os.Getenv("PVE_API_URL"),
		VMNames:     parseVMMap(os.Getenv("VM_NAMES")),
		HAProxyPath: os.Getenv("HAPROXY_PATH"),
		RalbUpdater: ralb,
		Logger:      logger,
		FetchDelay:  delay,
	}
}
