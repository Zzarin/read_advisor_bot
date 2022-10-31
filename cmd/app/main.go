package main

import (
	"context"
	"github.com/jessevdk/go-flags"
	"log"
	"read_advisor_bot/internal/api/client/telegram"
	"read_advisor_bot/internal/consumer"
	eventTelegram "read_advisor_bot/internal/events/telegram"
	"read_advisor_bot/internal/sqlite"
)

func main() {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal("Error parse env variables", err)
	}

	tgClient := telegram.New(cfg.TgBotHost, cfg.TgBotToken)
	storage, err := sqlite.NewDb(cfg.SqliteStoragePath)
	if err != nil {
		log.Fatalf("can't create new Db %v", err)
	}
	err = storage.Init(context.TODO())
	if err != nil {
		log.Fatalf("can't init Db %v", err)
	}
	eventsProcessor := eventTelegram.NewProcessor(tgClient, storage)
	log.Print("service started")
	handler := consumer.NewConsumer(eventsProcessor, eventsProcessor, cfg.BatchSize)
	err = handler.Start()
	if err != nil {
		log.Fatal("service is stopped", err)
	}
}
