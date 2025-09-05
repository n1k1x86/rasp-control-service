package server

import (
	"context"
	"log"
	"net/http"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"
	httphandlers "rasp-central-service/services/server/http_handlers"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	r   *mux.Router
	ctx context.Context
}

func (h *HTTPServer) Start() {
	log.Println("starting http server")
	http.ListenAndServe("0.0.0.0:8000", h.r)
}

func NewHTTPServer(ctx context.Context, ssrfRepo *ssrfrepo.Repository) *HTTPServer {
	r := httphandlers.BuildRouter(ssrfRepo)

	return &HTTPServer{
		r:   r,
		ctx: ctx,
	}
}
