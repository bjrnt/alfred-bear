package main

import (
	"fmt"
	"os"

	"github.com/deanishe/awgo"
	humanize "github.com/dustin/go-humanize"
)

var (
	IconWorkflow = &aw.Icon{"icon.png", aw.IconTypeImage}
	wf           *aw.Workflow
)

func init() {
	wf = aw.New()
}

func run() {
	query := os.Args[1]
	notes, _ := Search(query)
	for _, note := range notes {
		wf.NewItem(note.Title).Subtitle(fmt.Sprintf("Last edited %s", humanize.Time(note.ModificationDate))).UID(note.Identifier).Arg(note.Identifier).Valid(true).Icon(IconWorkflow)
	}
	wf.WarnEmpty("No matching notes found", "Try another query")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
