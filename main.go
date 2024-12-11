package main

import (
	"flag"
	"os"

	"google.golang.org/grpc"

	pb "github.com/IskanderSh/tages-task"
)

func main() {
	conn, err := grpc.Dial("localhost:1111", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pb.NewFileServiceClient(conn)

	var fileName string
	flag.StringVar(&fileName, "file-name", "test.txt", "file name")
	flag.Parse()

	_, err = os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	_, err = os.Open(fileName)
	if err != nil {
		panic(err)
	}
}
