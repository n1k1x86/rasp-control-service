package ssrfrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SSRFAgent struct {
	ID                 primitive.ObjectID `bson:"_id"`
	ServiceName        string             `bson:"service_name"`
	ServiceDescription string             `bson:"service_description"`
	IsActive           bool               `bson:"is_active"`
	Rules              Rules              `bson:"rules"`
	UpdateRulesURLs    string             `bson:"update_rules_url"`
}

type Rules struct {
	ID          primitive.ObjectID `bson:"_id"`
	HostRules   HostRules          `bson:"host_rules"`
	IPRules     IPRules            `bson:"ip_rules"`
	RegexpRules RegexpRules        `bson:"regexp_rules"`
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
