package general

import (
	"context"
	"errors"
	"fmt"
	"time"

	rasp_coll "rasp-central-service/services/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db  *mongo.Database
	ctx context.Context
}

func (r *Repository) RegService(serviceName, serviceDescription string) (string, error) {
	coll := r.db.Collection(rasp_coll.Services)

	service := CreateNewService(serviceName, serviceDescription)

	res, err := coll.InsertOne(r.ctx, service)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("inserted result is nil")
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *Repository) IsServiceRegistered(serviceID string) (bool, error) {
	coll := r.db.Collection(rasp_coll.Services)

	id, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id": id,
	}

	res := coll.FindOne(r.ctx, filter)
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) GetAllAgents() ([]*BaseAgent, error) {
	var agents []*BaseAgent = make([]*BaseAgent, 0)

	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{}

		cur, err := collection.Find(r.ctx, filter)
		if err != nil {
			return nil, err
		}
		results := make([]*BaseAgent, 0)
		err = cur.All(r.ctx, results)
		if err != nil {
			return nil, err
		}
		agents = append(agents, results...)
	}

	return agents, nil
}

func (r *Repository) GetAgentsByServiceName(serviceName string) ([]*BaseAgent, error) {
	var agents []*BaseAgent = make([]*BaseAgent, 0)

	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{
			"service_name": serviceName,
		}

		cur, err := collection.Find(r.ctx, filter)
		if err != nil {
			return nil, err
		}
		results := make([]*BaseAgent, 0)
		err = cur.All(r.ctx, results)
		if err != nil {
			return nil, err
		}
		agents = append(agents, results...)
	}

	return agents, nil
}

func (r *Repository) GetActiveAgentsCount() (int64, error) {
	var count int64
	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{
			"is_active": true,
		}

		res, err := collection.CountDocuments(r.ctx, filter)
		if err != nil {
			return 0, err
		}
		count += res
	}

	return count, nil
}

func (r *Repository) GetActiveAgentsCountByServiceName(serviceName string) (int64, error) {
	var count int64
	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{
			"service_name": serviceName,
			"is_active":    true,
		}

		res, err := collection.CountDocuments(r.ctx, filter)
		if err != nil {
			return 0, err
		}
		count += res
	}

	return count, nil
}

func (r *Repository) DeactivateAllAgents() (bool, error) {
	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{
			"is_active": true,
		}

		update := bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		}

		_, err := collection.UpdateMany(r.ctx, filter, update)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *Repository) DeactivateAllAgentsByServiceName(serviceName string) (bool, error) {
	for _, coll := range rasp_coll.CollectionsArray {
		collection := r.db.Collection(coll)

		filter := bson.M{
			"is_active": true,
		}

		update := bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		}

		_, err := collection.UpdateMany(r.ctx, filter, update)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func NewRepository(db *mongo.Database, ctx context.Context) *Repository {
	return &Repository{
		db:  db,
		ctx: ctx,
	}
}
