package sqlite

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
)

type Storage struct {
	db *sql.DB
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	_, err := io.WriteString(h, p.URL)
	if err != nil {
		return "", fmt.Errorf("can't calculate hash %w", err)
	}

	_, err = io.WriteString(h, p.UserName)
	if err != nil {
		return "", fmt.Errorf("can't calculate hash %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("can't create a database table %w", err)
	}
	return nil
}

//NewDb creates new SQLite storage.
func NewDb(path string) (*Storage, error) { //path is a path where Db is going to be stored
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't connect to database %w", err)
	}

	return &Storage{db: db}, nil
}

//Save saves a page.
func (s *Storage) Save(ctx context.Context, page *Page) (err error) {
	query := `INSERT INTO pages (url, user_name) VALUES (?, ?)`
	_, err = s.db.ExecContext(ctx, query, page.URL, page.UserName)
	if err != nil {
		return fmt.Errorf("can't save a page %w", err)
	}
	return nil
}

//PickRandom return a random page from previously saved pages.
func (s *Storage) PickRandom(ctx context.Context, userName string) (page *Page, err error) {
	query := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`
	var url string
	err = s.db.QueryRowContext(ctx, query, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't get a page %w", err)
	}
	return &Page{
		URL:      url,
		UserName: userName,
	}, nil
}

//Remove removes a page from the storage.
func (s *Storage) Remove(ctx context.Context, page *Page) error {
	query := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	_, err := s.db.ExecContext(ctx, query, page.URL, page.UserName)
	if err != nil {
		return fmt.Errorf("can't remove a page %w", err)
	}
	return nil
}

//IsExist checks if page exists in storage.
func (s *Storage) IsExist(ctx context.Context, page *Page) (bool, error) {
	query := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`
	var count int
	err := s.db.QueryRowContext(ctx, query, page.URL, page.UserName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("can't check if page exist %w", err)
	}

	return count > 0, nil
}
