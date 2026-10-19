package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

type Store struct {
	db       *sql.DB
	mu       sync.Mutex
	sessions map[string]string // token -> lowercase email, in-memory (nicht persistiert)
}

func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		email TEXT NOT NULL UNIQUE COLLATE NOCASE,
		password_hash TEXT NOT NULL
	)`); err != nil {
		return nil, err
	}
	if err := migrateVerifiedColumn(db); err != nil {
		return nil, err
	}
	if err := createStudioTables(db); err != nil {
		return nil, err
	}
	for _, table := range []string{verificationTable, resetTable} {
		if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS ` + table + ` (
			token_hash TEXT PRIMARY KEY,
			email TEXT NOT NULL COLLATE NOCASE,
			expires_at INTEGER NOT NULL
		)`); err != nil {
			return nil, err
		}
	}
	return &Store{db: db, sessions: make(map[string]string)}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// migrateVerifiedColumn fuegt users.verified hinzu, falls es fehlt. Konten, die schon vor
// der E-Mail-Bestaetigung existierten, werden dabei einmalig als bestaetigt markiert.
func migrateVerifiedColumn(db *sql.DB) error {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name = 'verified'").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	if _, err := db.Exec("ALTER TABLE users ADD COLUMN verified INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	_, err := db.Exec("UPDATE users SET verified = 1")
	return err
}
