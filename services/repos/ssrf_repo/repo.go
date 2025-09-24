package ssrfrepo

import (
	"context"
	"fmt"

	rasp_coll "rasp-central-service/services/database/mongo"
	generalrepo "rasp-central-service/services/repos/general"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db          *mongo.Database
	ctx         context.Context
	generalRepo *generalrepo.Repository
}

func (r *Repository) RegAgent(agent *SSRFAgent) (string, error) {
	isRegistered, err := r.generalRepo.IsServiceRegistered(agent.ServiceID.Hex())
	if err != nil {
		return "", err
	}
	if !isRegistered {
		return "", fmt.Errorf("service is not registered with id = %s", agent.ServiceID.Hex())
	}

	coll := r.db.Collection(rasp_coll.SSRFAgentsColl)
	res, err := coll.InsertOne(r.ctx, agent)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
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

func NewRepository(db *mongo.Database, ctx context.Context, generalRepo *generalrepo.Repository) *Repository {
	return &Repository{
		db:          db,
		ctx:         ctx,
		generalRepo: generalRepo,
	}
}
