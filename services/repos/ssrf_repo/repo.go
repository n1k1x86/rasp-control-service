package ssrfrepo

import (
	"context"
	"errors"
	"time"

	rasp_coll "rasp-central-service/services/database/mongo"

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
	coll := r.db.Collection("ssrf_agents")
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

	coll := r.db.Collection(rasp_coll.SSRFAgentsColl)
	_, err = coll.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateSSRFRules(agentID string, rules *Rules) error {
	coll := r.db.Collection(rasp_coll.SSRFAgentsColl)
	id, err := primitive.ObjectIDFromHex(agentID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"rules": rules,
		},
	}

	_, err = coll.UpdateByID(r.ctx, id, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteAgent(agentID string) error {
	coll := r.db.Collection(rasp_coll.SSRFAgentsColl)
	id, err := primitive.ObjectIDFromHex(agentID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": id,
	}

	_, err = coll.DeleteOne(r.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllAgents() ([]*SSRFAgent, error) {
	agents := make([]*SSRFAgent, 0)

	coll := r.db.Collection(rasp_coll.SSRFAgentsColl)
	filter := bson.M{}

	cur, err := coll.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cur.All(r.ctx, &agents)
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func NewRepository(db *mongo.Database, ctx context.Context) *Repository {
	return &Repository{
		db:  db,
		ctx: ctx,
	}
}
