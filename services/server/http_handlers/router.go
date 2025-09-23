package httphandlers

import (
	generalRepo "rasp-central-service/services/repos/general"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	"github.com/gorilla/mux"
	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

func BuildRouter(ssrfRepo *ssrfrepo.Repository, generalRepo *generalRepo.Repository, streams map[string]rasp_rpc.RASPCentral_SyncRulesServer) *mux.Router {
	r := mux.NewRouter()
	RegSSRFHandlers(r, ssrfRepo, streams)
	RegGeneralHandlers(r, generalRepo)
	return r
}
