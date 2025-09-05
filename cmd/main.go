package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"rasp-central-service/config"
	"rasp-central-service/services/database/mongo"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"
	"rasp-central-service/services/server"
	rasp_server "rasp-central-service/services/server"
	"syscall"
	"time"

	rasp_rpc "github.com/n1k1x86/rasp-grpc-contract/gen/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client, err := mongo.MongoConnect(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("mongo disconnected")
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	ssrfRepo := ssrfrepo.NewRepository(client.Database(cfg.Mongo.DBName), ctx)

	grpcServer := grpc.NewServer()
	server := server.NewGRPCServer(ctx, ssrfRepo)

	go func() {
		rasp_rpc.RegisterRASPCentralServer(grpcServer, server)
		log.Println("server is running")
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server error: %s", err.Error())
		}
	}()

	go func() {
		httpServer := rasp_server.NewHTTPServer(ctx, ssrfRepo)
		httpServer.Start()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	cancel()
	time.Sleep(cfg.App.GracefulTimeout)

	grpcServer.GracefulStop()
	log.Println("app was gracefully shutted down")
}
