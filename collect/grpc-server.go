package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	cpb "./proto"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Collect(ctx context.Context, req *cpb.CollectRequest) (*cpb.CollectResponse, error) {
	seed := fmt.Sprintf("https://%s", req.GetSeed())

	c := newCrawler(seed, nil)
	collected, err := c.crawl()
	if err != nil {
		return nil, err
	}

	return &cpb.CollectResponse{Collected: collected}, nil
}

func gRPCServerInit() {
	log.Println("Initializing gRPC server ...")

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

	log.Printf("gRPC server running at %s:%d", SHOST, SPORT)
}
