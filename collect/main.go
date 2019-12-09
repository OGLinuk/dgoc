package main

import (
	"log"

	cpb "github.com/OGLinuk/dgoc/collect/proto"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	collectClient cpb.CollectServiceClient
)

func init() {
	gRPCServerInit()
}

func main() {
	conn, err := grpc.Dial("0.0.0.0:8001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to grpc.Dial: %s", err.Error())
	}

	collectClient = cpb.NewCollectServiceClient(conn)

	log.Println("Successfully connected to collect gRPC server ...")

	grouter := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	grouter.GET("/:seed", CollectHandler)

	grouter.Use(cors.New(corsConfig))

	log.Println("Serving on 0.0.0.0:9001 ...")
	log.Fatal(grouter.Run("0.0.0.0:9001"), nil)
}
