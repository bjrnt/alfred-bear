package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type searchResponse struct {
	Notes string `json:"notes"`
}

type note struct {
	CreationDate     time.Time `json:"creationDate"`
	Title            string    `json:"title"`
	ModificationDate time.Time `json:"modificationDate"`
	Identifier       string    `json:"identifier"`
	Pin              string    `json:"pin"`
}

func formatSearch(query string) string {
	return fmt.Sprintf("\"bear://x-callback-url/search?show_window=no&term=%s&token=%s\"", query, bearToken)
}

func runCommand(query string) (searchResponse, error) {
	cmd := exec.Command(xcallPath, "-url", formatSearch(query))
	output, err := cmd.Output()
	fmt.Printf("%s\n", output)
	var response searchResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(output, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func search(query string) ([]note, error) {
	var notes []note
	result, err := runCommand(query)
	if err != nil {
		return notes, err
	}
	err = json.Unmarshal([]byte(result.Notes), &notes)
	return notes, err
}
