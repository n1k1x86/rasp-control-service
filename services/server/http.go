package server

import (
	"context"
	"log"
	"net/http"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"
	httphandlers "rasp-central-service/services/server/http_handlers"

	"github.com/gorilla/mux"
	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

type HTTPServer struct {
	r       *mux.Router
	ctx     context.Context
	streams map[string]rasp_rpc.RASPCentral_SyncRulesServer
}

func (h *HTTPServer) Start() {
	log.Println("starting http server")
	http.ListenAndServe("0.0.0.0:8000", h.r)
}

func NewHTTPServer(ctx context.Context, streams map[string]rasp_rpc.RASPCentral_SyncRulesServer, ssrfRepo *ssrfrepo.Repository) *HTTPServer {
	r := httphandlers.BuildRouter(ssrfRepo, streams)

	return &HTTPServer{
		r:       r,
		ctx:     ctx,
		streams: streams,
	}
}
