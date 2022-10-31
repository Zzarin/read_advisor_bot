package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"read_advisor_bot/internal/sqlite"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command `%s` from `%s`", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatID, text, userName)
	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)

	}
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	page := &sqlite.Page{
		URL:      pageURL,
		UserName: userName,
	}
	isExist, err := p.storage.IsExist(context.TODO(), page)
	if err != nil {
		return fmt.Errorf("can't check if page exist %w", err)
	}
	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	err = p.storage.Save(context.TODO(), page)
	if err != nil {
		return fmt.Errorf("can't save new page %w", err)
	}

	err = p.tg.SendMessage(chatID, msgSaved)
	if err != nil {
		return fmt.Errorf("can't send `Save` message to user %w", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {
	page, err := p.storage.PickRandom(context.TODO(), userName)
	if err != nil && !errors.Is(err, sqlite.ErrNoSavedPages) {
		return fmt.Errorf("can't pick a random page %w", err)
	}
	if errors.Is(err, sqlite.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	err = p.tg.SendMessage(chatID, page.URL)
	if err != nil {
		return fmt.Errorf("can't send page link %w", err)
	}

	return p.storage.Remove(context.TODO(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	if err == nil && u.Host != "" {
		return true
	}
	return false
}
