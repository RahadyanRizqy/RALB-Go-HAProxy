package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client *http.Client

func InitHTTPClient() {
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
}

func FetchVMs(cfg RalbEnv) ([]VM, error) {
	req, err := http.NewRequest("GET", cfg.PveAPIURL+"/api2/json/cluster/resources?type=vm", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", cfg.APIToken)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v\nraw: %s", err, string(body))
	}

	return result.Data, nil
}
