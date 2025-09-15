package ssrfrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SSRFAgent struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	AgentName          string             `json:"agent_name" bson:"agent_name"`
	ServiceName        string             `json:"service_name" bson:"service_name"`
	ServiceDescription string             `json:"service_description" bson:"service_description"`
	IsActive           bool               `json:"is_active" bson:"is_active"`
	Rules              Rules              `json:"rules" bson:"rules"`
	UpdateRulesURL     string             `json:"update_rules_url" bson:"update_rules_url"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
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
