package general

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseAgent struct {
	ID        primitive.ObjectID `bson:"_id"`
	AgentName string             `bson:"agent_name"`
	IsActive  bool               `bson:"is_active"`
	CreatedAt time.Duration      `bson:"created_at"`
	UpdatedAt time.Duration      `bson:"updated_at"`
}
