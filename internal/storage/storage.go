package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(page *Page) (err error)
	PickRandom(userName string) (page *Page, err error)
	Remove(page *Page) error
	IsExist(page *Page) (bool, error)
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
