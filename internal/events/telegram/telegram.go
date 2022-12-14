package telegram

import (
	"context"
	"errors"
	"fmt"
	tgClient "read_advisor_bot/internal/api/client/telegram"
	"read_advisor_bot/internal/events"
	"read_advisor_bot/internal/sqlite"
)

type Processor struct {
	tg      *tgClient.Client
	offset  int
	storage Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

type Storage interface {
	Save(ctx context.Context, page *sqlite.Page) (err error)
	PickRandom(ctx context.Context, userName string) (page *sqlite.Page, err error)
	Remove(ctx context.Context, page *sqlite.Page) error
	IsExist(ctx context.Context, page *sqlite.Page) (bool, error)
}

func NewProcessor(client *tgClient.Client, storage Storage) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit) //получаем апдейты
	if err != nil {
		return nil, fmt.Errorf("can't get event %w", err)
	}

	if len(updates) == 0 { //если список апдейтов пуст выходим
		return nil, nil
	}

	resp := make([]events.Event, 0, len(updates)) //готовим переменную для результата

	for _, u := range updates { //перебираем все апдейты и преобразуем в тип эвент
		resp = append(resp, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 //обновляем оффсет. чтобы получать свежие апдейты
	//апдейты относятся только к телеге и в других мессанджерах их возможно не будет
	//эвент это структура нашего приложения к которой мы преобразуем сообщения от разных
	//мессенджеров в каком бы формате они не приходили
	return resp, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("unknown event type")
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't get event %w", err)
	}

	err = p.doCmd(event.Text, meta.ChatID, meta.UserName)
	if err != nil {
		return fmt.Errorf("can't perform a command %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("can't get Meta")
	}
	return res, nil
}

func event(update tgClient.Update) events.Event {
	updateType := fetchType(update)
	resp := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}
	if updateType == events.Message {
		resp.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.Username,
		}
	}
	return resp
}

func fetchText(update tgClient.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update tgClient.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
