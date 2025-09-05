package httphandlers

import ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

type GetAllAgentsResponse struct {
	Agents []*ssrfrepo.SSRFAgent `json:"agents"`
}
