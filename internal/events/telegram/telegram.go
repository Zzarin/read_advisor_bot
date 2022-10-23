package telegram

import "read_advisor_bot/internal/api/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	//storage
}

func New(client *telegram.Client) {

}
