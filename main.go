package main

import (
	"fmt"
	"os"

	aw "github.com/deanishe/awgo"
	humanize "github.com/dustin/go-humanize"
	"golang.org/x/text/unicode/norm"
)

var (
	IconWorkflow = &aw.Icon{
		Value: "icon.png",
		Type:  aw.IconTypeImage,
	}
	wf          *aw.Workflow
	userHomeDir string
)

func init() {
	wf = aw.New()

	// Try to read the user's home directory. This is required for finding Bear's SQLite DB
	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		wf.FatalError(err)
	}
}

func run() {
	// For some annoying reason, UTF-8 arguments from Alfred are not normalized in the correct way,
	// so any non-ASCII characters break the query. Here we normalize the UTF-8 input to make sure
	// that it can be handled properly by the DB
	query := string(norm.NFC.Bytes([]byte(os.Args[1])))

	// search the user's bear notes
	notes, err := search(query)
	if err != nil {
		wf.FatalError(err)
	}

	// send the results to alfred
	sendNotes(notes)
	wf.SendFeedback()
}

func sendNotes(notes []note) {
	for _, note := range notes {
		wf.NewItem(note.Title).
			Subtitle(fmt.Sprintf("Last edited %s", humanize.Time(note.ModificationDate))).
			UID(note.Identifier).
			Arg(note.Identifier).
			Valid(true).
			Icon(IconWorkflow)
	}
	wf.WarnEmpty("No matching notes found", "Try another query")
}

func main() {
	wf.Run(run)
}
