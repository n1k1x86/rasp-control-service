package ssrfrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SSRFAgent struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	AgentName string             `json:"agent_name" bson:"agent_name"`
	ServiceID primitive.ObjectID `json:"service_id" bson:"service_id"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	Rules     Rules              `json:"rules" bson:"rules"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Rules struct {
	HostRules   HostRules   `json:"host_rules" bson:"host_rules"`
	IPRules     IPRules     `json:"ip_rules" bson:"ip_rules"`
	RegexpRules RegexpRules `json:"regexp_rules" bson:"regexp_rules"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at"`
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

func (r *Repository) NewRules(hosts, ips, regexps []string) *Rules {
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
		HostRules:   hostRules,
		IPRules:     ipRules,
		RegexpRules: regRules,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (r *Repository) NewAgent(agentName, serviceID string) (*SSRFAgent, error) {
	rules := &Rules{}
	serviceIDObjectID, err := primitive.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, err
	}
	return &SSRFAgent{
		ID:        primitive.NewObjectID(),
		AgentName: agentName,
		ServiceID: serviceIDObjectID,
		IsActive:  true,
		Rules:     *rules,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
