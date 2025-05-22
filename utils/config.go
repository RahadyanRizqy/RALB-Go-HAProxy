package utils

import (
	"fmt"
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

	// DEFAULT VALUES
	ralb, err := strconv.ParseBool(os.Getenv("RALB_UPDATER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		ralb = false
	}

	delay, err := strconv.Atoi(os.Getenv("FETCH_DELAY"))
	if err != nil {
		delay = 1000
	}

	logger, err := strconv.ParseBool(os.Getenv("LOGGER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		ralb = false
	}

	runserver, err := strconv.ParseBool(os.Getenv("RUN_SERVER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		runserver = false
	}

	return RalbEnv{
		APIToken:    os.Getenv("API_TOKEN"),
		PveAPIURL:   os.Getenv("PVE_API_URL"),
		VMNames:     parseVMMap(os.Getenv("VM_NAMES")),
		HAProxyPath: os.Getenv("HAPROXY_PATH"),
		RalbUpdater: ralb,
		Logger:      logger,
		RunServer:   runserver,
		FetchDelay:  delay,
	}
}
