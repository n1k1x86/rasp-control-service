package httphandlers

import ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

type ErrorResponse struct {
	Detail string `json:"detail"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type GetAllAgentsResponse struct {
	Agents []*ssrfrepo.SSRFAgent `json:"agents"`
}

type BindRulesRequest struct {
	AgentID string `json:"agent_id"`
	RulesID string `json:"rules_id"`
}

type AddNewSSRFRulesRequest struct {
	Description string   `json:"description"`
	HostsRules  []string `json:"hosts_rules"`
	IPRules     []string `json:"ip_rules"`
	RegexpRules []string `json:"regexp_rules"`
}

type AddNewSSRFRulesResponse struct {
	RulesID string `json:"rules_id"`
	Detail  string `json:"detail"`
}

type RegServiceResponse struct {
	ServiceID string `json:"service_id"`
	Detail    string `json:"detail"`
}

type RegServiceRequest struct {
	ServiceName        string `json:"service_name"`
	ServiceDescription string `json:"service_description"`
}
