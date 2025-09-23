package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"rasp-central-service/config"
	"rasp-central-service/services/database/mongo"
	generalrepo "rasp-central-service/services/repos/general"
	ssrfrepo "rasp-central-service/services/repos/ssrf_repo"
	rasp_server "rasp-central-service/services/server"
	"syscall"

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

	generalRepo := generalrepo.NewRepository(client.Database(cfg.Mongo.DBName), ctx)
	ssrfRepo := ssrfrepo.NewRepository(client.Database(cfg.Mongo.DBName), ctx, generalRepo)

	grpcServer := grpc.NewServer()
	server := rasp_server.NewGRPCServer(ctx, ssrfRepo, generalRepo)

	go func() {
		rasp_rpc.RegisterRASPCentralServer(grpcServer, server)
		log.Println("server is running")
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server error: %s", err.Error())
		}
	}()

	go func() {
		httpServer := rasp_server.NewHTTPServer(ctx, server.StreamMap, ssrfRepo, generalRepo)
		httpServer.Start()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	cancel()

	gracefulCtx, gracefulCancel := context.WithTimeout(context.Background(), cfg.App.GracefulTimeout)
	defer gracefulCancel()
	<-gracefulCtx.Done()

	grpcServer.GracefulStop()
	log.Println("app was gracefully shutted down")
}
