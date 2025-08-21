package ssrfrepo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db  *mongo.Database
	ctx context.Context
}

func NewRepository(db *mongo.Database, ctx context.Context) *Repository {
	return &Repository{
		db:  db,
		ctx: ctx,
	}
}
