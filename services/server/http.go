package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	generalRepo "rasp-central-service/services/repos/general"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"
	httphandlers "rasp-central-service/services/server/http_handlers"
	"time"

	"github.com/gorilla/mux"
	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
)

type HTTPServer struct {
	r       *mux.Router
	ctx     context.Context
	streams map[string]rasp_rpc.RASPCentral_SyncRulesServer
	server  *http.Server
}

func (h *HTTPServer) Start() error {
	log.Println("starting http server")
	err := h.server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("error while starting http server: %s", err.Error())
	}
	return nil
}

func (h *HTTPServer) Close(ctx context.Context) error {
	err := h.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	log.Println("http was stopped")
	return nil
}

func NewHTTPServer(ctx context.Context, addr string, streams map[string]rasp_rpc.RASPCentral_SyncRulesServer, ssrfRepo *ssrfrepo.Repository, generalRepo *generalRepo.Repository) *HTTPServer {
	r := httphandlers.BuildRouter(ssrfRepo, generalRepo, streams)

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,

		// Ограничение количества одновременных подключений
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	return &HTTPServer{
		r:       r,
		ctx:     ctx,
		streams: streams,
		server:  server,
	}
}
