package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc-stream/chat"
	"log"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func handleRecv(wg *sync.WaitGroup, client chat.GreeterClient, stream chat.Greeter_SayHelloClient) {
	defer wg.Done()

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		//print message to console
		fmt.Printf("\n \t\t Server : %s \n", resp.Reply)

	}

}
func handleSend(wg *sync.WaitGroup, client chat.GreeterClient, stream chat.Greeter_SayHelloClient) {
	defer wg.Done()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Client :")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failure in Reading string")
		}
		str = strings.Trim(str, "\r\n")

		err = stream.Send(&chat.HelloRequest{Greeting: str})
		if err != nil {
			log.Fatalf("failure in Reading string")
		}
		fmt.Printf("\n")

	}
}

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Unable to dial %v", err)
	}
	defer conn.Close()
	fmt.Println("Client connection Established")
	client := chat.NewGreeterClient(conn)
	stream, err := client.SayHello(context.Background())
	if err != nil {
		log.Fatalf("failure in SayHello")
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	defer wg.Wait()
	//ctx := context.Background()
	go handleSend(wg, client, stream)
	go handleRecv(wg, client, stream)

	//handleStream(client, stream)

}
