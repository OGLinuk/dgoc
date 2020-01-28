package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	qpb "github.com/OGLinuk/dgoc/queue/proto"
	"google.golang.org/grpc"
)

/*

// https://gist.github.com/moraes/2141121

type Queue []*Node

func (q *Queue) Push(n *Node) {
    *q = append(*q, n)
}

func (q *Queue) Pop() (n *Node) {
    n = (*q)[0]
    *q = (*q)[1:]
    return
}

func (q *Queue) Len() int {
    return len(*q)
}

*/

type server struct {
	//rc *RedisClient
	Queue    []string
	dequeued map[string]struct{} // For duplicate checking when queuing
}

func (s *server) Size(ctx context.Context, req *qpb.QueueSizeRequest) (*qpb.QueueSizeResponse, error) {
	log.Printf("%v", req)

	return &qpb.QueueSizeResponse{Size: int64(len(s.Queue)), Queued: s.Queue}, nil
}

func (s *server) Push(ctx context.Context, req *qpb.QueuePushRequest) (*qpb.QueuePushResponse, error) {
	enqueue := req.GetEnqueue()

	if _, exists := s.dequeued[enqueue]; !exists {
		s.Queue = append(s.Queue, enqueue)
		s.dequeued[enqueue] = struct{}{}
	}

	resp := &qpb.QueuePushResponse{
		Success: true,
	}

	return resp, nil
}

func (s *server) Pop(ctx context.Context, req *qpb.QueuePopRequest) (*qpb.QueuePopResponse, error) {
	//k := req.GetKey() // Really only for API auth
	v := s.Queue[0]
	s.Queue = s.Queue[1:]

	return &qpb.QueuePopResponse{Dequeued: v}, nil
}

func gRPCServerInit() {
	log.Println("Initializing gRPC server ...")

	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "8003")

	SHOST := os.Getenv("SERVER_HOST")
	SPORT, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Failed strconv.Atoi: %s", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", SHOST, SPORT))
	if err != nil {
		log.Fatalf("Failed to net.Listen: %s", err.Error())
	}

	srvr := &server{
		//rc: NewRedisClient(),
		Queue: make([]string, 1),
	}

	/*
		if err = srvr.rc.Init(); err != nil {
			log.Fatalf("Failed to srvr.rc.Init: %s", err.Error())
		}
	*/

	srv := grpc.NewServer()
	qpb.RegisterQueueServiceServer(srv, srvr)

	log.Printf("gRPC server running at %s:%d", SHOST, SPORT)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
