package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()

	//processor = processor.New()

	//consumer.Start(fetcher, processor)
}

func mustToken() string {
	//bot -tg-bot-token `my token`
	token := flag.String(
		"token-bot-token",
		"",
		"token for access telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not passed")
	}

	return *token
}
