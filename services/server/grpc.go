package server

import (
	"context"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

type RASPCentral struct {
	rasp_rpc.UnimplementedRASPCentralServer
	ctx       context.Context
	SSRFRepo  *ssrfrepo.Repository
	StreamMap map[string]rasp_rpc.RASPCentral_SyncRulesServer
}

func (r *RASPCentral) RegSSRFAgent(ctx context.Context, req *rasp_rpc.RegSSRFAgentRequest) (*rasp_rpc.RegSSRFAgentResponse, error) {
	rules := r.SSRFRepo.NewRules(req.HostRules, req.IPRules, req.RegexpRules)
	agent := r.SSRFRepo.NewAgent(rules, req.ServiceName, req.ServiceDescription, req.UpdateURL, req.AgentName)
	id, err := r.SSRFRepo.RegAgent(agent)
	if err != nil {
		return nil, err
	}

	resp := &rasp_rpc.RegSSRFAgentResponse{
		Status:  200,
		Detail:  "agent was sucessfully registered",
		AgentID: id,
	}

	return resp, nil
}

func (r *RASPCentral) CloseSSRFAgent(ctx context.Context, req *rasp_rpc.AgentRequest) (*rasp_rpc.CloseSSRFAgentResponse, error) {
	err := r.SSRFRepo.DeleteAgent(req.AgentID)
	if err != nil {
		return nil, err
	}
	return &rasp_rpc.CloseSSRFAgentResponse{Detail: "agent was close successfuly"}, nil
}

func (r *RASPCentral) DeactivateSSRFAgent(ctx context.Context, req *rasp_rpc.DeactivateSSRFAgentRequest) (*rasp_rpc.DeactivateSSRFAgentResponse, error) {
	err := r.SSRFRepo.DeactivateAgent(req.AgentID)
	if err != nil {
		return nil, err
	}

	resp := &rasp_rpc.DeactivateSSRFAgentResponse{
		Status: 200,
		Detail: "agent " + req.AgentName + " was deactivated",
	}

	return resp, nil
}

func (r *RASPCentral) SyncRules(req *rasp_rpc.AgentRequest, stream rasp_rpc.RASPCentral_SyncRulesServer) error {
	r.StreamMap[req.AgentID] = stream
	<-r.ctx.Done()
	return nil
}

func NewGRPCServer(ctx context.Context, ssrfRepo *ssrfrepo.Repository) *RASPCentral {
	return &RASPCentral{
		SSRFRepo:  ssrfRepo,
		ctx:       ctx,
		StreamMap: make(map[string]rasp_rpc.RASPCentral_SyncRulesServer),
	}
}
