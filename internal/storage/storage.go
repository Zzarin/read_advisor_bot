package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

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
