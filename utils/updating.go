package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func ModifyHAProxy(cfg RalbEnv, vmResults []VMMetric) {
	data, err := os.ReadFile(cfg.HAProxyPath)
	if err != nil {
		log.Printf("Failed to read HAProxy config: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")

	// Locate the backend section and the server lines
	backendStart := -1
	checkLine := -1
	for i, line := range lines {
		if backendStart == -1 && strings.HasPrefix(strings.TrimSpace(line), "backend web_servers") {
			backendStart = i
		}
		if backendStart != -1 && strings.Contains(strings.TrimSpace(line), "option httpchk GET /course") {
			checkLine = i
			break
		}
	}
	if backendStart == -1 || checkLine == -1 {
		log.Println("Unable to locate 'backend web_servers' or 'option httpchk' line")
		return
	}

	// Collect all existing server lines and their indentation
	serverLines := make(map[string]string)
	var indent string
	serverStart := checkLine + 1
	serverEnd := serverStart
	for i := serverStart; i < len(lines); i++ {
		trim := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trim, "server") {
			fields := strings.Fields(trim)
			if len(fields) >= 2 {
				serverLines[fields[1]] = lines[i]
				if indent == "" {
					indent = lines[i][:strings.Index(lines[i], "server")]
				}
				serverEnd = i + 1
			}
		} else {
			break // stop at the first non-server line
		}
	}

	// Reorder server lines based on vmResults
	var reordered []string
	for _, vm := range vmResults {
		if line, ok := serverLines[vm.Name]; ok {
			reordered = append(reordered, line)
		}
	}

	// Replace the old server lines with reordered ones
	newLines := append([]string{}, lines[:serverStart]...)
	newLines = append(newLines, reordered...)
	newLines = append(newLines, lines[serverEnd:]...)

	// Save the result back to file
	if err := os.WriteFile(cfg.HAProxyPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		log.Printf("Failed to write updated HAProxy config: %v\n", err)
		fmt.Println("| Status : ERROR 				   |")
		fmt.Println("+----------------------------------------------+")
		fmt.Printf("Error Messsage : %v\n", err)
	} else {
		fmt.Println("| Status : SUKSES			       |")
		fmt.Println("+----------------------------------------------+")
	}
}
