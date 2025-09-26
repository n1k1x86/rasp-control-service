package httphandlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"

	"github.com/gorilla/mux"
)

func RegSSRFHandlers(r *mux.Router, ssrfRepo *ssrfrepo.Repository, streams map[string]rasp_rpc.RASPCentral_SyncRulesServer) {
	baseURI := "/ssrf-agents"

	bindRulesURL := baseURI + "/bind-rules"
	getAllURL := baseURI + "/get-all"
	addRulesURL := baseURI + "/add-rules"

	r.HandleFunc(getAllURL, GetAllSSRFAgents(ssrfRepo)).Methods("GET")
	logHandlers(getAllURL, "get")

	r.HandleFunc(bindRulesURL, BindSSRFRules(streams, ssrfRepo)).Methods("POST")
	logHandlers(bindRulesURL, "post")

	r.HandleFunc(addRulesURL, AddNewSSRFRules(ssrfRepo)).Methods("POST")
	logHandlers(addRulesURL, "post")
}

func GetAllSSRFAgents(ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	agents, err := ssrfRepo.GetAllAgents()
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(&agents)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func BindSSRFRules(streams map[string]rasp_rpc.RASPCentral_SyncRulesServer, ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			HandleError(w, err, http.StatusBadRequest)
			return
		}
		var req BindRulesRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}

		rules, err := ssrfRepo.BindRules(req.AgentID, req.RulesID)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}
		log.Printf("agent %s, rules were successfully updated", req.AgentID)

		stream := streams[req.AgentID]

		err = stream.Send(BuildRules(rules.HostRules.Hosts, rules.IPRules.IPs, rules.RegexpRules.Regexps))
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
	}
}

func AddNewSSRFRules(ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			HandleError(w, err, http.StatusBadRequest)
			return
		}

		var req AddNewSSRFRulesRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}

		rules := ssrfRepo.NewRules(req.IPRules, req.HostsRules, req.RegexpRules, req.Description)
		rulesID, err := ssrfRepo.AddRules(rules)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}

		resp := AddNewSSRFRulesResponse{RulesID: rulesID, Detail: "new rules were added"}
		data, err = json.Marshal(&resp)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
