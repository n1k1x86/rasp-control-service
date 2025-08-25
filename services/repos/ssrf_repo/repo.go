package ssrfrepo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db  *mongo.Database
	ctx context.Context
}

func (r *Repository) IsActivated(coll *mongo.Collection, agentName string) (bool, string, error) {
	filter := bson.M{
		"agent_name": agentName,
	}

	res := coll.FindOne(r.ctx, filter)
	var agent SSRFAgent
	err := res.Decode(&agent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, "", nil
		}
		return false, "", err
	}
	if agent.IsActive {
		return true, agent.ID.Hex(), nil
	}
	return false, "", nil
}

func (r *Repository) RegAgent(agent *SSRFAgent) (string, error) {
	coll := r.db.Collection(collName)

	ok, id, err := r.IsActivated(coll, agent.AgentName)
	if err != nil {
		return "", err
	}
	if ok {
		return id, nil
	}

	filter := bson.M{
		"agent_name": agent.AgentName,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  true,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"_id":              agent.ID,
			"agent_name":       agent.AgentName,
			"service_name":     agent.ServiceName,
			"rules":            agent.Rules,
			"update_rules_url": agent.UpdateRulesURL,
			"created_at":       agent.CreatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	res := coll.FindOneAndUpdate(r.ctx, filter, update, opts)

	var doc SSRFAgent
	err = res.Decode(&doc)
	if err != nil {
		return "", err
	}

	return doc.ID.Hex(), nil
}

func (r *Repository) DeactivateAgent(agentID string) error {
	id, err := primitive.ObjectIDFromHex(agentID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	coll := r.db.Collection(collName)
	_, err = coll.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *mongo.Database, ctx context.Context) *Repository {
	return &Repository{
		db:  db,
		ctx: ctx,
	}
}
