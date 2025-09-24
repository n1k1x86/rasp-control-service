package httphandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"

	"github.com/gorilla/mux"
)

func RegSSRFHandlers(r *mux.Router, ssrfRepo *ssrfrepo.Repository, streams map[string]rasp_rpc.RASPCentral_SyncRulesServer) {
	baseURI := "/ssrf-agents"
	updateURL := baseURI + "/update"
	getAllURL := baseURI + "/get-all"

	r.HandleFunc(getAllURL, GetAllSSRFAgents(ssrfRepo)).Methods("GET")
	logHandlers(getAllURL, "get")

	r.HandleFunc(updateURL, UpdateSSRFRules(streams, ssrfRepo)).Methods("POST")
	logHandlers(updateURL, "post")
}

func GetAllSSRFAgents(ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	agents, err := ssrfRepo.GetAllAgents()
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": %s}`, err.Error())))
			return
		}

		data, err := json.Marshal(&agents)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": %s}`, err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func UpdateSSRFRules(streams map[string]rasp_rpc.RASPCentral_SyncRulesServer, ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var req UpdateRulesBody
		err = json.Unmarshal(data, &req)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rules := ssrfRepo.NewRules(req.HostsRules, req.IPRules, req.RegexpRules)
		err = ssrfRepo.UpdateSSRFRules(req.AgentID, rules)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("agent %s, rules were successfully updated", req.AgentID)

		stream := streams[req.AgentID]

		err = stream.Send(BuildRules(req.HostsRules, req.IPRules, req.RegexpRules))
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
	}
}
