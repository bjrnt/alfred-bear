package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// note is a Bear note
type note struct {
	Identifier       string // The identifier is an UUID that can be used to open this note in Bear if selected
	Title            string
	CreationDate     time.Time
	ModificationDate time.Time
}

const (
	// See: https://bear.app/faq/Where%20are%20Bear%27s%20notes%20located/
	BEAR_DB_LOC = "/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

	// The primary query to search for notes in Bear's SQLite database. Bear uses Core Data under
	// the hood, which is why table and column names are a bit strange. Core data encodes dates as
	// floats with a different base date from UNIX timestamps, so some date fiddling is required
	// (see: https://stackoverflow.com/a/2923127/300664). The query turns them back into regular
	// UNIX timestamps. Bear sets a ZTRASHEDDATE for notes that have been deleted, so we exclude
	// those.
	NOTE_QUERY = `
	SELECT 
		ZUNIQUEIDENTIFIER, 
		ZTITLE, 
		strftime('%s', ZCREATIONDATE, 'unixepoch', '31 years'),
		strftime('%s', ZMODIFICATIONDATE, 'unixepoch', '31 years')
	FROM ZSFNOTE 
	WHERE 
		ZTITLE LIKE ? OR ZTEXT LIKE ? 
		AND ZTRASHEDDATE IS NULL 
	ORDER BY ZMODIFICATIONDATE DESC`
)

// search the user's notes for the given query.
func search(query string) ([]note, error) {
	db, err := openDB()
	if err != nil {
		return nil, errors.Wrap(err, "open sqlite db")
	}
	defer db.Close()

	// col LIKE '%abc%' will search for abc anywhere in the VARCHAR column, ignoring case. Reformat
	// the query to be surrounded by wildcard % chars.
	query = fmt.Sprintf("%%%s%%", query)
	// execute query with escaped parameters
	rows, err := db.Query(NOTE_QUERY, query, query)
	if err != nil {
		return nil, errors.Wrap(err, "execute sqlite query")
	}

	notes, err := parseNotes(rows)
	if err != nil {
		return nil, errors.Wrap(err, "parse sqlite rows")
	}
	return notes, nil
}

// openDB opens Bear's SQLite database in readonly mode
func openDB() (*sql.DB, error) {
	dbLoc := fmt.Sprintf("file:%s%s?mode=ro", userHomeDir, BEAR_DB_LOC)
	return sql.Open("sqlite3", dbLoc)
}

// parseNotes from SQL query result rows
func parseNotes(rows *sql.Rows) ([]note, error) {
	defer rows.Close()

	notes := []note{}

	for rows.Next() {
		note := note{}

		var creationDate int64
		var modificationDate int64

		// Columns need to be scanned in the order they are selected in the NOTE_QUERY
		err := rows.Scan(&note.Identifier, &note.Title, &creationDate, &modificationDate)
		if err != nil {
			return nil, errors.Wrap(err, "scan note DB column")
		}

		note.CreationDate = time.Unix(creationDate, 0)
		note.ModificationDate = time.Unix(modificationDate, 0)

		notes = append(notes, note)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return notes, nil
}
