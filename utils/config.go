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

	logger, err := strconv.ParseBool(os.Getenv("LOGGER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		logger = false
	}

	fetchDelay, err := strconv.Atoi(os.Getenv("FETCH_DELAY"))
	if err != nil {
		fetchDelay = 1000
	}

	netIfaceRate, err := strconv.Atoi(os.Getenv("NET_IFACE_RATE"))
	if err != nil {
		netIfaceRate = 12500000
	}

	serverStart, err := strconv.ParseBool(os.Getenv("RUN_SERVER"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		serverStart = false
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		serverPort = 9000
	}

	serverSuccessMessage := os.Getenv("SERVER_SUCCESS_MESSAGE")
	if serverSuccessMessage == "" {
		serverSuccessMessage = "RALB OK!"
	}

	serverErrorMessage := os.Getenv("SERVER_ERROR_MESSAGE")
	if serverErrorMessage == "" {
		serverErrorMessage = "RALB NOT OK!"
	}

	return RalbEnv{
		APIToken:             os.Getenv("API_TOKEN"),
		PveAPIURL:            os.Getenv("PVE_API_URL"),
		VMNames:              parseVMMap(os.Getenv("VM_NAMES")),
		HAProxyPath:          os.Getenv("HAPROXY_PATH"),
		RalbUpdater:          ralbUpdater,
		Logger:               logger,
		FetchDelay:           fetchDelay,
		NetIfaceRate:         netIfaceRate,
		ServerStart:          serverStart,
		ServerPort:           serverPort,
		ServerSuccessMessage: serverSuccessMessage,
		ServerErrorMessage:   serverErrorMessage,
	}
}
