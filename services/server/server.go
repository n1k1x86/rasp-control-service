package server

import (
	"context"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RASPCentral struct {
	Server   *rasp_rpc.UnimplementedRASPCentralServer
	SSRFRepo *ssrfrepo.Repository
}

func (r *RASPCentral) RegAgent(ctx context.Context, req *rasp_rpc.RegAgentRequest) (*rasp_rpc.RegAgentResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method RegAgent not implemented")
}

func (r *RASPCentral) DeactivateAgent(ctx context.Context, req *rasp_rpc.DeactivateAgentRequest) (*rasp_rpc.DeactivateAgentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateAgent not implemented")
}

func NewServer(ctx context.Context, ssrfRepo *ssrfrepo.Repository) *RASPCentral {
	return &RASPCentral{
		Server:   &rasp_rpc.UnimplementedRASPCentralServer{},
		SSRFRepo: ssrfRepo,
	}
}
