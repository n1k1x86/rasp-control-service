package httphandlers

import (
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"

	"github.com/gorilla/mux"
)

func BuildRouter(ssrfRepo *ssrfrepo.Repository) *mux.Router {
	r := mux.NewRouter()
	RegSSRFHandlers(r, ssrfRepo)
	return r
}
