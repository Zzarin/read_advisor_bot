package main

import (
	"flag"
	"log"
	"read_advisor_bot/internal/api/client/telegram"
	"read_advisor_bot/internal/consumer"
	eventTelegram "read_advisor_bot/internal/events/telegram"
	"read_advisor_bot/internal/storage/files"
)

func main() {
	tgClient := telegram.New("api.telegram.org", mustToken())
	storage := files.NewStorage("files_storage")
	eventsProcessor := eventTelegram.NewProcessor(tgClient, storage)
	log.Print("service started")
	handler := consumer.NewConsumer(eventsProcessor, eventsProcessor, 100)
	err := handler.Start()
	if err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	//bot -tg-bot-token `my token`
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not passed")
	}

	return *token
}
