package storage

import (
	"database/sql"

	// register sqlite driver
	_ "modernc.org/sqlite"
)

// Store implements db storage
type Store struct {
	db *sql.DB
}

// NewStore instantiates the storage
func NewStore(dsn string) *Store {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	// Create the table if it does not exist
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS number (
			id INTEGER PRIMARY KEY,
			cc NCHAR(3) NOT NULL,
			ndc NCHAR(4) NOT NULL,
			sn NCHAR(13) NOT NULL, 
			used BOOLEAN NOT NULL DEFAULT 0, 
			domain TEXT NOT NULL,
			carrier TEXT NOT NULL,
			userID  INTEGER NOT NULL DEFAULT 0, 
			allocated INTEGER NOT NULL DEFAULT 0, 
			reserved  INTEGER NOT NULL DEFAULT 0, 
			deallocated INTEGER NOT NULL DEFAULT 0, 
			portedIn  INTEGER NOT NULL DEFAULT 0, 
			portedOut INTEGER NOT NULL DEFAULT 0, 
			CONSTRAINT unq UNIQUE (cc, ndc, sn)
		);
		`); err != nil {
		panic(err)
	}
	return &Store{db: db}
}

//Close closes db connection
func (s *Store) Close() {
	s.db.Close()
}
