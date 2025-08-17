package mongo

import (
	"context"
	"log"
	"rasp-central-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(ctx context.Context, cfg config.Mongo) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(cfg.ApplyURI).
		SetMaxPoolSize(cfg.MaxConns).
		SetMinPoolSize(cfg.MinConns).
		SetMaxConnIdleTime(cfg.IdleTime))
	if err != nil {
		return nil, err
	}
	log.Println("connected to mongo")
	return client, nil
}
