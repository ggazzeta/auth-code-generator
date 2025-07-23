package store

import (
	"auth-code-generator/pkg/models"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3" // The sqlite3 driver
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(dataSourceName string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	query := `
    CREATE TABLE IF NOT EXISTS codes (
        user_id TEXT PRIMARY KEY,
        user_email TEXT NOT NULL,
        code TEXT NOT NULL,
        created_at DATETIME NOT NULL
    );`
	if _, err = db.Exec(query); err != nil {
		return nil, err
	}

	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) Save(code models.StoredCode) error {
	query := `
    INSERT INTO codes (user_id, user_email, code, created_at)
    VALUES (?, ?, ?, ?)
    ON CONFLICT(user_id) DO UPDATE SET
        user_email = excluded.user_email,
        code = excluded.code,
        created_at = excluded.created_at;`

	_, err := s.db.Exec(query, code.UserID, code.UserEmail, code.Code, code.CreatedAt)
	return err
}

func (s *SqliteStore) Get(userID string) (models.StoredCode, bool, error) {
	var storedCode models.StoredCode
	var createdAt string // Read time as string to parse it correctly

	query := "SELECT user_id, user_email, code, created_at FROM codes WHERE user_id = ?;"
	row := s.db.QueryRow(query, userID)

	err := row.Scan(&storedCode.UserID, &storedCode.UserEmail, &storedCode.Code, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.StoredCode{}, false, nil
		}
		return models.StoredCode{}, false, err
	}

	storedCode.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return models.StoredCode{}, false, err
	}

	return storedCode, true, nil
}
