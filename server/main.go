package main

import (
	"bufio"
	"fmt"
	"grpc-stream/chat"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type GreeterServer struct {
	chat.UnimplementedGreeterServer
}

func handleRecv(wg *sync.WaitGroup, stream chat.Greeter_SayHelloServer) error {
	defer wg.Done()

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			log.Println("EOF Recieved ")
		}
		if err != nil {
			log.Fatalf("Error Recieving message %v ", err)
		}
		if len(req.GetGreeting()) > 0 {
			fmt.Println("\n \t\t Client :", req.GetGreeting())
		}

	}
}

func handleSend(wg *sync.WaitGroup, stream chat.Greeter_SayHelloServer) error {
	defer wg.Done()

	for {
		fmt.Printf("Server :")
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error Reading message  ", err)
		}
		str = strings.Trim(str, "\r\n")
		err = stream.Send(&chat.HelloResponse{Reply: str})
		if err != nil {
			fmt.Printf("Failed to send Message %v \n", err)
			log.Fatalf("Error Sending message %v ", err)
			return err
		}
		fmt.Printf("\n")
	}
}

func (g *GreeterServer) SayHello(stream chat.Greeter_SayHelloServer) error {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go handleRecv(wg, stream)
	go handleSend(wg, stream)

	wg.Wait()
	return nil

}

func main() {
	lis, err := net.Listen("tcp", "localhost:5000")

	if err != nil {
		log.Fatalf("Error in Listen %v", err)
	}
	fmt.Printf("Listening on localhost:5000 \n")

	defer lis.Close()
	grpcServer := grpc.NewServer()

	chat.RegisterGreeterServer(grpcServer, &GreeterServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Unable to serve GRPC : %v", err)
	}

}
