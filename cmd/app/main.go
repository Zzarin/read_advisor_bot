package main

import (
	"context"
	"flag"
	"log"
	"read_advisor_bot/internal/api/client/telegram"
	"read_advisor_bot/internal/consumer"
	eventTelegram "read_advisor_bot/internal/events/telegram"
	"read_advisor_bot/internal/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = ".data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
	storage, err := sqlite.NewDb(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't create new Db %v", err)
	}
	err = storage.Init(context.TODO())
	if err != nil {
		log.Fatalf("can't init Db %v", err)
	}
	eventsProcessor := eventTelegram.NewProcessor(tgClient, storage)
	log.Print("service started")
	handler := consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)
	err = handler.Start()
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
