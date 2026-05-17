package db

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

var instance *sql.DB

func Init(path string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Panic("Failed to create db directory:", err)
	}

	var err error
	instance, err = sql.Open("sqlite", path)
	if err != nil {
		log.Panic("Failed to open SQLite database:", err)
	}

	_, err = instance.Exec(`CREATE TABLE IF NOT EXISTS test_notes (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		text       TEXT    NOT NULL,
		created_at TEXT    NOT NULL
	)`)
	if err != nil {
		log.Panic("Failed to init test_notes table:", err)
	}

	log.Println("SQLite initialized:", path)
}

func InsertNote(text string) error {
	_, err := instance.Exec(
		`INSERT INTO test_notes (text, created_at) VALUES (?, ?)`,
		text, time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}

type Note struct {
	ID        int
	Text      string
	CreatedAt string
}

func GetAllNotes() ([]Note, error) {
	rows, err := instance.Query(`SELECT id, text, created_at FROM test_notes ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Text, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}
