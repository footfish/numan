package storage

import (
	"database/sql"

	"github.com/footfish/numan"
	// register sqlite driver
	_ "modernc.org/sqlite"
)

// store implements db storage
type store struct {
	db *sql.DB
}

// NewStore instantiates the storage
func NewStore(dsn string) numan.API {
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
	return &store{db: db}
}

//Close closes db connection
func (s *store) Close() {
	s.db.Close()
}
