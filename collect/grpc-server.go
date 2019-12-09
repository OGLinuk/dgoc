package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	cpb "github.com/OGLinuk/dgoc/collect/proto"
	spb "github.com/OGLinuk/dgoc/store/proto"
	"google.golang.org/grpc"
)

type server struct{}

var (
	storeClient spb.StoreServiceClient
)

func (s *server) Collect(ctx context.Context, req *cpb.CollectRequest) (*cpb.CollectResponse, error) {
	seed := fmt.Sprintf("https://%s", req.GetSeed())

	c := newCrawler(seed, nil)
	collected, err := c.crawl()
	if err != nil {
		return nil, err
	}

	sReq := &spb.StoreRequest{
		Crawled:   seed,
		Collected: collected,
	}

	_, err = storeClient.Store(ctx, sReq)
	if err != nil {
		return nil, err
	}

	return &cpb.CollectResponse{Collected: collected}, nil
}

func gRPCStoreInit() {
	log.Println("Initializing connection to store gRPC server ...")

	conn, err := grpc.Dial("store-service:8002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to grpc.Dial: %s", err.Error())
	}

	storeClient = spb.NewStoreServiceClient(conn)

	log.Println("Successfully connected to store gRPC server ...")
}

func gRPCServerInit() {
	gRPCStoreInit()

	log.Println("Initializing collect gRPC server ...")

	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "8001")

	SHOST := os.Getenv("SERVER_HOST")
	SPORT, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Failed strconv.Atoi: %s", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", SHOST, SPORT))
	if err != nil {
		log.Fatalf("Failed to net.Listen: %s", err.Error())
	}

	srvr := grpc.NewServer()
	cpb.RegisterCollectServiceServer(srvr, &server{})

	go func() {
		if err := srvr.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %s", err.Error())
		}
	}()

	log.Printf("collect gRPC server running at %s:%d ...", SHOST, SPORT)
}
