package httphandlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	generalRepo "rasp-central-service/services/repos/general"

	"github.com/gorilla/mux"
)

func RegGeneralHandlers(r *mux.Router, generalRepo *generalRepo.Repository) {
	baseURI := "/general"
	regServiceURL := baseURI + "/reg-service"
	healthURL := baseURI + "/health"

	r.HandleFunc(regServiceURL, RegService(generalRepo)).Methods("POST")
	logHandlers(regServiceURL, "post")

	r.HandleFunc(healthURL, Health()).Methods("GET")
	logHandlers(healthURL, "get")
}

func RegService(generalRepo *generalRepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req RegServiceRequest

		err = json.Unmarshal(data, &req)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := generalRepo.RegService(req.ServiceName, req.ServiceDescription)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := RegServiceResponse{ServiceID: id, Detail: "service was successfully registered"}

		body, err := json.Marshal(&resp)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)
		w.WriteHeader(http.StatusOK)
	}
}

func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{Status: "OK"}
		data, err := json.Marshal(&resp)
		if err != nil {
			HandleError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
