package server

import (
	"context"
	"log"
	generalRepo "rasp-central-service/services/repos/general"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

type RASPCentral struct {
	rasp_rpc.UnimplementedRASPCentralServer
	ctx         context.Context
	SSRFRepo    *ssrfrepo.Repository
	GeneralRepo *generalRepo.Repository
	StreamMap   map[string]rasp_rpc.RASPCentral_SyncRulesServer
}

func (r *RASPCentral) RegSSRFAgent(ctx context.Context, req *rasp_rpc.RegSSRFAgentRequest) (*rasp_rpc.RegSSRFAgentResponse, error) {
	agent, err := r.SSRFRepo.NewAgent(req.AgentName, req.ServiceID)
	if err != nil {
		return nil, err
	}
	id, err := r.SSRFRepo.RegAgent(agent)
	if err != nil {
		return nil, err
	}

	resp := &rasp_rpc.RegSSRFAgentResponse{
		Status:  200,
		Detail:  "agent was sucessfully registered",
		AgentID: id,
	}
	log.Println("agent was registered with id =", id)

	return resp, nil
}

func (r *RASPCentral) CloseSSRFAgent(ctx context.Context, req *rasp_rpc.AgentRequest) (*rasp_rpc.CloseSSRFAgentResponse, error) {
	err := r.SSRFRepo.DeleteAgent(req.AgentID)
	if err != nil {
		return nil, err
	}
	log.Println("agent was closed with id=", req.AgentID)
	return &rasp_rpc.CloseSSRFAgentResponse{Detail: "agent was close successfuly"}, nil
}

func (r *RASPCentral) IsServiceRegistered(ctx context.Context, req *rasp_rpc.IsServiceRegisteredReq) (*rasp_rpc.IsServiceRegisteredResp, error) {
	isRegistered, err := r.GeneralRepo.IsServiceRegistered(req.GetServiceID())
	if err != nil {
		return nil, err
	}

	return &rasp_rpc.IsServiceRegisteredResp{IsRegistered: isRegistered}, nil
}

func (r *RASPCentral) SyncRules(req *rasp_rpc.AgentRequest, stream rasp_rpc.RASPCentral_SyncRulesServer) error {
	r.StreamMap[req.AgentID] = stream
	log.Println("synced updater connection with agent ", req.AgentID)
	<-r.ctx.Done()
	return nil
}

func NewGRPCServer(ctx context.Context, ssrfRepo *ssrfrepo.Repository, generalRepo *generalRepo.Repository) *RASPCentral {
	return &RASPCentral{
		SSRFRepo:    ssrfRepo,
		GeneralRepo: generalRepo,
		ctx:         ctx,
		StreamMap:   make(map[string]rasp_rpc.RASPCentral_SyncRulesServer),
	}
}
