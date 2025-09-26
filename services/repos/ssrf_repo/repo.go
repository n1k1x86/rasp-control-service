package ssrfrepo

import (
	"context"
	"errors"
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

var ErrNoAgent = fmt.Errorf("agent does not exist")
var ErrNoRules = fmt.Errorf("rules do not exist")

func (r *Repository) IsAgentExist(agentID primitive.ObjectID, coll *mongo.Collection) (bool, error) {
	filter := bson.M{
		"_id": agentID,
	}

	res := coll.FindOne(r.ctx, filter)
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
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

func (r *Repository) AddRules(rules *Rules) (string, error) {
	coll := r.db.Collection(rasp_coll.SSRFRulesColl)

	res, err := coll.InsertOne(r.ctx, rules)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *Repository) DeleteRules(rulesID string) (bool, error) {
	coll := r.db.Collection(rasp_coll.SSRFRulesColl)

	objID, err := primitive.ObjectIDFromHex(rulesID)
	if err != nil {
		return false, err
	}

	ok, err := r.IsRulesExist(objID, coll)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("rules do not exist")
	}

	filter := bson.M{
		"_id": objID,
	}

	_, err = coll.DeleteOne(r.ctx, filter)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) GetRulesByID(rulesID primitive.ObjectID, coll *mongo.Collection) (*Rules, error) {
	filter := bson.M{
		"_id": rulesID,
	}
	res := coll.FindOne(r.ctx, filter)
	err := res.Err()
	if err != nil {
		return nil, err
	}

	var rules Rules
	err = res.Decode(&rules)
	if err != nil {
		return nil, err
	}
	return &rules, nil
}

func (r *Repository) IsRulesExist(rulesID primitive.ObjectID, coll *mongo.Collection) (bool, error) {
	filter := bson.M{"_id": rulesID}
	res := coll.FindOne(r.ctx, filter)
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) BindRules(agentID, rulesID string) (*Rules, error) {
	agentObjID, err := primitive.ObjectIDFromHex(agentID)
	if err != nil {
		return nil, err
	}

	agentColl := r.db.Collection(rasp_coll.SSRFAgentsColl)
	ok, err := r.IsAgentExist(agentObjID, agentColl)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoAgent
	}

	rulesObjID, err := primitive.ObjectIDFromHex(rulesID)
	if err != nil {
		return nil, err
	}

	rulesColl := r.db.Collection(rasp_coll.SSRFRulesColl)
	ok, err = r.IsRulesExist(rulesObjID, rulesColl)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoRules
	}

	rules, err := r.GetRulesByID(rulesObjID, rulesColl)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"rules_id": rulesObjID,
		},
	}

	_, err = agentColl.UpdateByID(r.ctx, agentObjID, update)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func NewRepository(db *mongo.Database, ctx context.Context, generalRepo *generalrepo.Repository) *Repository {
	return &Repository{
		db:          db,
		ctx:         ctx,
		generalRepo: generalRepo,
	}
}
