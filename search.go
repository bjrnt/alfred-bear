package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// searchResponse holds the structure returned by Bear from the x-callback
type searchResponse struct {
	Notes *string `json:"notes"`
}

// note is a Bear note
type note struct {
	CreationDate     time.Time `json:"creationDate"`
	Title            string    `json:"title"`
	ModificationDate time.Time `json:"modificationDate"`
	Identifier       string    `json:"identifier"`
	Pin              string    `json:"pin"`
}

func runSearch(query string) ([]byte, error) {
	callback := fmt.Sprintf("\"bear://x-callback-url/search?show_window=no&term=%s&token=%s\"", url.PathEscape(query), bearToken)
	log.Printf(url.PathEscape(query))
	cmd := exec.Command(xcallPath, "-url", callback)
	return cmd.Output()
}

func parseResponse(searchOutput []byte) ([]note, error) {
	var response searchResponse
	if err := json.Unmarshal(searchOutput, &response); err != nil {
		return []note{}, errors.Wrap(err, "could not unmarshal search response")
	}
	// empty response
	if response.Notes == nil {
		return []note{}, nil
	}
	// non-empty response
	var notes []note
	err := json.Unmarshal([]byte(*response.Notes), &notes)
	return notes, errors.Wrap(err, "could not unmarshal notes")
}

func search(query string) ([]note, error) {
	output, err := runSearch(query)
	if err != nil {
		return []note{}, errors.Wrap(err, "could not perform search")
	}
	log.Printf("%s\n", output)
	notes, err := parseResponse(output)
	return notes, errors.Wrap(err, "could not parse search result")
}
