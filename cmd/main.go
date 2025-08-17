package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"rasp-central-service/config"
	"rasp-central-service/services/database/mongo"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.MongoConnect(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("mongo disconnected")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	cancel()
}
