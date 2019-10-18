package main

import (
	"encoding/json"
	"fmt"
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

const searchUrl = `"bear://x-callback-url/search?show_window=no&term=%s&token=%s"`

func runSearch(query string) ([]byte, error) {
	cbUrl := fmt.Sprintf(searchUrl, url.PathEscape(query), bearToken)
	cmd := exec.Command(xcallPath, "-url", cbUrl)
	return cmd.Output()
}

func parseResponse(searchOutput []byte) ([]note, error) {
	var response searchResponse
	err := json.Unmarshal(searchOutput, &response)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal search response")
	}
	// empty response
	if response.Notes == nil {
		return nil, nil
	}
	// non-empty response
	var notes []note
	err = json.Unmarshal([]byte(*response.Notes), &notes)
	return notes, errors.Wrap(err, "could not unmarshal notes")
}

func search(query string) ([]note, error) {
	output, err := runSearch(query)
	if err != nil {
		return nil, errors.Wrap(err, "could not perform search")
	}
	notes, err := parseResponse(output)
	return notes, errors.Wrap(err, "could not parse search result")
}
