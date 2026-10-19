package store

import (
	"database/sql"
	"errors"
)

// ── LAB-STUDIO ────────────────────────────────────────────────────────────────────────────
// Jede Abfrage enthaelt "owner = ?": ein Konto sieht und aendert nur seine eigenen Daten.

var ErrNotFound = errors.New("nicht gefunden")

type Song struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type View struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Scale float64 `json:"scale"`
}

type Field struct {
	ID     int64   `json:"id"`
	Type   string  `json:"type"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	W      float64 `json:"w"`
	H      float64 `json:"h"`
	Text   string  `json:"text"`
	SongID int64   `json:"song_id"`
}

func createStudioTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			owner TEXT NOT NULL COLLATE NOCASE,
			title TEXT NOT NULL,
			text TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS studio_fields (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			owner TEXT NOT NULL COLLATE NOCASE,
			type TEXT NOT NULL,
			x REAL NOT NULL,
			y REAL NOT NULL,
			w REAL NOT NULL,
			h REAL NOT NULL,
			text TEXT NOT NULL DEFAULT '',
			song_id INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS studio_view (
			owner TEXT PRIMARY KEY COLLATE NOCASE,
			x REAL NOT NULL,
			y REAL NOT NULL,
			scale REAL NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListSongs(owner string) ([]Song, error) {
	rows, err := s.db.Query("SELECT id, title FROM songs WHERE owner = ? ORDER BY title COLLATE NOCASE", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := []Song{}
	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.ID, &song.Title); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, rows.Err()
}

func (s *Store) CreateSong(owner, title string) (Song, error) {
	result, err := s.db.Exec("INSERT INTO songs (owner, title) VALUES (?, ?)", owner, title)
	if err != nil {
		return Song{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Song{}, err
	}
	return Song{ID: id, Title: title}, nil
}

func (s *Store) GetSong(owner string, id int64) (Song, error) {
	var song Song
	err := s.db.QueryRow("SELECT id, title, text FROM songs WHERE id = ? AND owner = ?", id, owner).
		Scan(&song.ID, &song.Title, &song.Text)
	if errors.Is(err, sql.ErrNoRows) {
		return Song{}, ErrNotFound
	}
	return song, err
}

func (s *Store) SaveSongText(owner string, id int64, text string) error {
	result, err := s.db.Exec("UPDATE songs SET text = ? WHERE id = ? AND owner = ?", text, id, owner)
	return requireOneRow(result, err)
}

// LoadView liefert die gespeicherte Ansicht oder die Startansicht (0, 0, 100 %).
func (s *Store) LoadView(owner string) (View, error) {
	view := View{Scale: 1}
	err := s.db.QueryRow("SELECT x, y, scale FROM studio_view WHERE owner = ?", owner).
		Scan(&view.X, &view.Y, &view.Scale)
	if errors.Is(err, sql.ErrNoRows) {
		return View{Scale: 1}, nil
	}
	return view, err
}

func (s *Store) SaveView(owner string, view View) error {
	_, err := s.db.Exec(`INSERT INTO studio_view (owner, x, y, scale) VALUES (?, ?, ?, ?)
		ON CONFLICT(owner) DO UPDATE SET x = excluded.x, y = excluded.y, scale = excluded.scale`,
		owner, view.X, view.Y, view.Scale)
	return err
}

const fieldColumns = "id, type, x, y, w, h, text, song_id"

func scanField(row interface{ Scan(...any) error }) (Field, error) {
	var f Field
	err := row.Scan(&f.ID, &f.Type, &f.X, &f.Y, &f.W, &f.H, &f.Text, &f.SongID)
	return f, err
}

func (s *Store) ListFields(owner string) ([]Field, error) {
	rows, err := s.db.Query("SELECT "+fieldColumns+" FROM studio_fields WHERE owner = ? ORDER BY id", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fields := []Field{}
	for rows.Next() {
		f, err := scanField(rows)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (s *Store) GetField(owner string, id int64) (Field, error) {
	row := s.db.QueryRow("SELECT "+fieldColumns+" FROM studio_fields WHERE id = ? AND owner = ?", id, owner)
	f, err := scanField(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Field{}, ErrNotFound
	}
	return f, err
}

func (s *Store) CreateField(owner string, f Field) (Field, error) {
	result, err := s.db.Exec(`INSERT INTO studio_fields (owner, type, x, y, w, h, text, song_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, owner, f.Type, f.X, f.Y, f.W, f.H, f.Text, f.SongID)
	if err != nil {
		return Field{}, err
	}
	f.ID, err = result.LastInsertId()
	return f, err
}

func (s *Store) UpdateField(owner string, f Field) error {
	result, err := s.db.Exec(`UPDATE studio_fields SET x = ?, y = ?, w = ?, h = ?, text = ?, song_id = ?
		WHERE id = ? AND owner = ?`, f.X, f.Y, f.W, f.H, f.Text, f.SongID, f.ID, owner)
	return requireOneRow(result, err)
}

// requireOneRow meldet ErrNotFound, wenn ein UPDATE keine Zeile getroffen hat
// (unbekannte ID oder fremdes Konto).
func requireOneRow(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
