package ssrfrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SSRFAgent struct {
	ID        primitive.ObjectID `bson:"_id"`
	AgentName string             `bson:"agent_name"`
	ServiceID primitive.ObjectID `bson:"service_id"`
	IsActive  bool               `bson:"is_active"`
	Rules     primitive.ObjectID `bson:"rules_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Rules struct {
	ID          primitive.ObjectID `bson:"_id"`
	Description string             `bson:"description"`
	HostRules   HostRules          `bson:"host_rules"`
	IPRules     IPRules            `bson:"ip_rules"`
	RegexpRules RegexpRules        `bson:"regexp_rules"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type HostRules struct {
	Hosts     []string  `json:"hosts" bson:"hosts"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type IPRules struct {
	IPs       []string  `json:"ips" bson:"ips"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type RegexpRules struct {
	Regexps   []string  `json:"regexps" bson:"regexps"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (r *Repository) NewRules(hosts, ips, regexps []string, description string) *Rules {
	hostRules := HostRules{
		Hosts:     hosts,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ipRules := IPRules{
		IPs:       ips,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	regRules := RegexpRules{
		Regexps:   regexps,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &Rules{
		ID:          primitive.NewObjectID(),
		Description: description,
		HostRules:   hostRules,
		IPRules:     ipRules,
		RegexpRules: regRules,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (r *Repository) NewAgent(agentName, serviceID string) (*SSRFAgent, error) {
	serviceIDObjectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, err
	}
	return &SSRFAgent{
		ID:        primitive.NewObjectID(),
		AgentName: agentName,
		ServiceID: serviceIDObjectID,
		IsActive:  true,
		Rules:     primitive.NilObjectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
