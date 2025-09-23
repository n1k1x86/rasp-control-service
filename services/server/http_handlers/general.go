package httphandlers

import (
	"encoding/json"
	"io"
	"net/http"
	generalRepo "rasp-central-service/services/repos/general"

	"github.com/gorilla/mux"
)

func RegGeneralHandlers(r *mux.Router, generalRepo *generalRepo.Repository) {
	baseURI := "/general"

	r.HandleFunc(baseURI+"/reg_service", RegService(generalRepo)).Methods("POST")
	logHandlers("registered route /general/reg_service with method POST")
}

func RegService(generalRepo *generalRepo.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			logHandlers(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req RegServiceRequest

		err = json.Unmarshal(data, &req)
		if err != nil {
			logHandlers(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := generalRepo.RegService(req.ServiceName, req.ServiceDescription)
		if err != nil {
			logHandlers(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := RegServiceResponse{ServiceID: id, Detail: "service was successfully registered"}

		body, err := json.Marshal(&resp)
		if err != nil {
			logHandlers(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)
		w.WriteHeader(http.StatusOK)
	}
}
