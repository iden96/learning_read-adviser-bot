package main

import (
	"context"
	"flag"
	tgClient "iden69/read-adviser-bot/clients/telegram"
	event_consumer "iden69/read-adviser-bot/consumer/event-consumer"
	"iden69/read-adviser-bot/events/telegram"
	"iden69/read-adviser-bot/storage/sqlite"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	filesStoragePath  = "files_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	t := mustToken()
	tgClient := tgClient.New(tgBotHost, t)
	// storage := files.New(storagePath)
	storage, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: ", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(tgClient, storage)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
