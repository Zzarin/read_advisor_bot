package main

type Config struct {
	TgBotToken        string `long:"tg-bot-token" description:"Token to start the bot" env:"TG_BOT_TOKEN"`
	TgBotHost         string `long:"tg-bot-host" description:"Bot host" env:"TG_BOT_HOST"`
	SqliteStoragePath string `long:"sqlite-storage-path" description:"Path to store sqlite db " env:"SQLITE_STORAGE_PATH"`
	BatchSize         int    `long:"batch-size" description:"How many updates to get from telegram API" env:"BATCH_SIZE"`
}
