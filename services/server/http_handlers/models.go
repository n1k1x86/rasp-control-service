package httphandlers

import ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

type GetAllAgentsResponse struct {
	Agents []*ssrfrepo.SSRFAgent `json:"agents"`
}

type UpdateRulesBody struct {
	AgentID     string   `json:"agent_id"`
	HostsRules  []string `json:"hosts_rules"`
	IPRules     []string `json:"ip_rules"`
	RegexpRules []string `json:"regexp_rules"`
}

type RegServiceResponse struct {
	ServiceID string `json:"service_id"`
	Detail    string `json:"detail"`
}

type RegServiceRequest struct {
	ServiceName        string `json:"service_name"`
	ServiceDescription string `json:"service_description"`
}
