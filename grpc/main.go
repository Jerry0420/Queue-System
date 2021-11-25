package main

import (
	"fmt"

	grpcServices "github.com/jerry0420/queue-system/grpc/proto"
)

func main() {
	fmt.Println(grpcServices.NewGrpcServiceClient)
}
