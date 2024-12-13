package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/IskanderSh/tages-task/proto"
)

var client pb.FileProviderClient
var fileName string

func main() {
	conn, err := grpc.Dial("localhost:1111", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = pb.NewFileProviderClient(conn)

	file := flag.String("file", "test.txt", "path to file")
	function := flag.String("func", "fetch", "which function to call")
	flag.Parse()

	fmt.Printf("running file-fetch-client with function: %s and using file: %s", *function, *file)
	fileName = *file

	switch *function {
	case "upload":
		upload()
	case "download":
		download()
	case "fetch":
		fetch()
	}
}

func upload() {
	_, err := os.Stat(fileName)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024*1024) // 1mb
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(string(buffer[:n]))

		// Отправляем чанк серверу
		err = stream.Send(&pb.UploadFileRequest{
			FileName: file.Name(),
			Content:  buffer[:n],
		})
		if err != nil {
			panic(err)
		}
	}

	// Завершаем загрузку
	resp, err := stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Имя загруженного файла: %v", resp.FileName)
}

func download() {
	stream, err := client.DownloadFile(context.Background(), &pb.DownloadFileRequest{FileName: fileName})
	if err != nil {
		panic(err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("finishing receiving chunks")
			break
		} else if err != nil {
			panic(err)
		}

		fmt.Printf("received new chunk: %s", string(response.Content))
	}
}

func fetch() {
	response, err := client.FetchFiles(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("fetch response data: %+v", response.Data)
}