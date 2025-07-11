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
	ralbUpdater, err := strconv.ParseBool(os.Getenv("RALB_UPDATER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		ralbUpdater = false
	}

	updateNotify, err := strconv.ParseBool(os.Getenv("UPDATE_NOTIFY"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		updateNotify = false
	}

	haproxyWeight, err := strconv.Atoi(os.Getenv("HAPROXY_WEIGHT"))
	if err != nil {
		fmt.Println("2. Error parsing boolean:", err)
		haproxyWeight = 256
	}

	consolePrint, err := strconv.ParseBool(os.Getenv("CONSOLE_PRINT"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		consolePrint = false
	}

	logger, err := strconv.ParseBool(os.Getenv("LOGGER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		logger = false
	}

	fetchDelay, err := strconv.Atoi(os.Getenv("FETCH_DELAY"))
	if err != nil {
		fetchDelay = 1000
	}

	netIfaceRate, err := strconv.ParseFloat(os.Getenv("NET_IFACE_RATE"), 64)
	if err != nil {
		netIfaceRate = 12500000
	}

	strict, err := strconv.ParseBool(os.Getenv("STRICT"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		strict = false
	}

	return RalbEnv{
		APIToken:       os.Getenv("API_TOKEN"),
		PveAPIURL:      os.Getenv("PVE_API_URL"),
		HAProxySock:    os.Getenv("HAPROXY_SOCK"),
		HAProxyBackend: os.Getenv("HAPROXY_BACKEND"),
		VMNames:        parseVMMap(os.Getenv("VM_NAMES")),
		VMIPs:          parseVMMap(os.Getenv("VM_IPS")),
		HAProxyPath:    os.Getenv("HAPROXY_PATH"),
		RalbUpdater:    ralbUpdater,
		UpdateNotify:   updateNotify,
		HAProxyWeight:  haproxyWeight,
		ConsolePrint:   consolePrint,
		Logger:         logger,
		FetchDelay:     fetchDelay,
		NetIfaceRate:   netIfaceRate,
		Strict:         strict,
	}
}
