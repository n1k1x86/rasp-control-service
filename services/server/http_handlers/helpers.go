package httphandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

func logHandlers(url, method string) {
	log.Printf("route %s was registered with method %s", url, strings.ToUpper(method))
}

func BuildErrorResponse(detail string) ([]byte, error) {
	resp := ErrorResponse{Detail: detail}
	data, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func HandleError(w http.ResponseWriter, err error, respStatus int) {
	log.Printf("ERROR: %s", err.Error())
	w.WriteHeader(respStatus)
	body, err := BuildErrorResponse(err.Error())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
}

func BuildSSRFRules(hosts, ips, regexps []string) *rasp_rpc.NewRules {
	payload := &rasp_rpc.NewRules_SSRFRules{
		SSRFRules: &rasp_rpc.UpdatedSSRFRules{
			HostRules:   hosts,
			IPRules:     ips,
			RegexpRules: regexps,
		},
	}
	return &rasp_rpc.NewRules{
		Payload: payload,
	}
}
