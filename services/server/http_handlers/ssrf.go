package httphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	"github.com/gorilla/mux"
)

func RegSSRFHandlers(r *mux.Router, ssrfRepo *ssrfrepo.Repository) {
	r.HandleFunc("/ssrf-agents/get_all", GetAllSSRFAgents(ssrfRepo)).Methods("GET")
}

func GetAllSSRFAgents(ssrfRepo *ssrfrepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	agents, err := ssrfRepo.GetAllAgents()
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": %s}`, err.Error())))
			return
		}

		data, err := json.Marshal(&agents)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": %s}`, err.Error())))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
