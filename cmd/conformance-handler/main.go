package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Initialize registry for object references
	registry := NewRegistry()
	defer registry.Cleanup()

	// Read requests from stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse request
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendResponse(NewHandlerErrorResponse("", "INVALID_REQUEST", ""))
			continue
		}

		resp := handleRequest(registry, req)
		sendResponse(resp)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}

// sendResponse writes a response to stdout as JSON
func sendResponse(resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}

	fmt.Println(string(data))
}
