package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	qpb "github.com/OGLinuk/dgoc/queue/proto"
	spb "github.com/OGLinuk/dgoc/store/proto"
	"google.golang.org/grpc"
)

type server struct{}

var (
	mdb     *MongoStore
	qClient qpb.QueueServiceClient
)

func (s *server) Store(ctx context.Context, req *spb.StoreRequest) (*spb.StoreResponse, error) {
	crawled := req.GetCrawled()
	parsed, err := url.Parse(crawled)
	if err != nil {
		return nil, err
	}

	if err = mdb.PutCrawled(parsed.Hostname(), []string{crawled}); err != nil {
		return nil, err
	}

	collected := req.GetCollected()

	if err = mdb.PutUncrawled("queued", collected); err != nil {
		return nil, err
	}

	for _, c := range collected {
		qClient.Push(ctx, &qpb.QueuePushRequest{Enqueue: c})
	}

	log.Printf("Collected: %d", len(collected))

	qSize, err := qClient.Size(ctx, &qpb.QueueSizeRequest{Key: "tmp"})
	if err != nil {
		return nil, err
	}

	log.Printf("Queue size: %d", qSize.Size)

	return &spb.StoreResponse{Success: true}, nil
}

func init() {
	mdb = NewMongoStore()

	if err := mdb.Init("mongodb"); err != nil {
		log.Fatalf("Failed to mdb.Init: %s", err.Error())
	}
}

func gRPCQueueInit() {
	log.Println("Initializing connection to queue gRPC server ...")

	conn, err := grpc.Dial("queue-service:8003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to grpc.Dial: %s", err.Error())
	}

	qClient = qpb.NewQueueServiceClient(conn)

	log.Println("Successfully connected to queue gRPC server ...")
}

func gRPCServerInit() {
	gRPCQueueInit()

	log.Println("Initializing store gRPC server ...")

	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "8002")

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
	spb.RegisterStoreServiceServer(srvr, &server{})

	log.Printf("store gRPC server running at %s:%d ...", SHOST, SPORT)

	if err := srvr.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
