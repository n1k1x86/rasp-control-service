package ssrfrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collName = "ssrf_agents"

type SSRFAgent struct {
	ID                 primitive.ObjectID `bson:"_id"`
	AgentName          string             `bson:"agent_name"`
	ServiceName        string             `bson:"service_name"`
	ServiceDescription string             `bson:"service_description"`
	IsActive           bool               `bson:"is_active"`
	Rules              Rules              `bson:"rules"`
	UpdateRulesURL     string             `bson:"update_rules_url"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

type Rules struct {
	ID          primitive.ObjectID `bson:"_id"`
	HostRules   HostRules          `bson:"host_rules"`
	IPRules     IPRules            `bson:"ip_rules"`
	RegexpRules RegexpRules        `bson:"regexp_rules"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type HostRules struct {
	ID        primitive.ObjectID `bson:"_id"`
	Hosts     []string           `bson:"hosts"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type IPRules struct {
	ID        primitive.ObjectID `bson:"_id"`
	IPs       []string           `bson:"ips"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type RegexpRules struct {
	ID        primitive.ObjectID `bson:"_id"`
	Regexps   []string           `bson:"regexps"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (r *Repository) NewRules(hosts, ips, regexps []string) *Rules {
	hostRules := HostRules{
		ID:        primitive.NewObjectID(),
		Hosts:     hosts,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ipRules := IPRules{
		ID:        primitive.NewObjectID(),
		IPs:       ips,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	regRules := RegexpRules{
		ID:        primitive.NewObjectID(),
		Regexps:   regexps,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &Rules{
		ID:          primitive.NewObjectID(),
		HostRules:   hostRules,
		IPRules:     ipRules,
		RegexpRules: regRules,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (r *Repository) NewAgent(rules *Rules, serviceName, description, updateURL, agentName string) *SSRFAgent {
	return &SSRFAgent{
		ID:                 primitive.NewObjectID(),
		AgentName:          agentName,
		ServiceName:        serviceName,
		ServiceDescription: description,
		IsActive:           true,
		Rules:              *rules,
		UpdateRulesURL:     updateURL,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
