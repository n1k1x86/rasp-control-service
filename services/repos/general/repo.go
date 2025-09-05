package general

import (
	"context"
	"time"

	rasp_coll "rasp-central-service/services/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db  *mongo.Database
	ctx context.Context
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
