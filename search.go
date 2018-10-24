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

type Note struct {
	CreationDate     time.Time `json:"creationDate"`
	Title            string    `json:"title"`
	ModificationDate time.Time `json:"modificationDate"`
	Identifier       string    `json:"identifier"`
	Pin              string    `json:"pin"`
}

const bearToken = "E822A5-0FA4D9-EE3A6D"
const xcallPath = "/Applications/xcall.app/Contents/MacOS/xcall"

func formatSearch(query string) string {
	return fmt.Sprintf("\"bear://x-callback-url/search?show_window=no&term=%s&token=%s\"", query, bearToken)
}

func runCommand(query string) (searchResponse, error) {
	cmd := exec.Command(xcallPath, "-url", formatSearch(query))
	output, err := cmd.Output()
	response := searchResponse{}
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(output, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func Search(query string) ([]Note, error) {
	notes := make([]Note, 0)
	result, err := runCommand(query)
	if err != nil {
		return notes, err
	}
	err = json.Unmarshal([]byte(result.Notes), &notes)
	return notes, err
}
