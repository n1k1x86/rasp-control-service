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

type Service struct {
	ID                 primitive.ObjectID `bson:"_id"`
	ServiceName        string             `bson:"service_name"`
	ServiceDescription string             `bson:"service_description"`
	RegDate            time.Time          `bson:"reg_date"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

func CreateNewService(serviceName, serviceDescription string) *Service {
	return &Service{
		ID:                 primitive.NewObjectID(),
		ServiceName:        serviceName,
		ServiceDescription: serviceDescription,
		RegDate:            time.Now(),
		UpdatedAt:          time.Now(),
	}
}
