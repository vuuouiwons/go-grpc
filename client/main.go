package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "go-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	if !scanner.Scan() {
		os.Exit(0)
	}
	return strings.Trim(scanner.Text(), "\r\n\t \ufeff\x00")
}

func readBool(scanner *bufio.Scanner, prompt string) bool {
	input := strings.ToLower(readLine(scanner, prompt))
	return input == "y" || input == "yes" || input == "true" || input == "1"
}

func callNormal(client pb.SimpleServiceClient, scanner *bufio.Scanner) {
	message := readLine(scanner, "Enter message: ")
	errorParam := readBool(scanner, "Trigger error? (y/n): ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Normal(ctx, &pb.Request{
		Message: message,
		Error:   errorParam,
	})
	if err != nil {
		fmt.Printf("Received error: %v\n", err)
		return
	}

	fmt.Printf("Received response: %s\n", response.GetMessage())
}

func callServerStream(client pb.SimpleServiceClient, scanner *bufio.Scanner) {
	message := readLine(scanner, "Enter message: ")
	errorParam := readBool(scanner, "Trigger error? (y/n): ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ServerStream(ctx, &pb.Request{
		Message: message,
		Error:   errorParam,
	})
	if err != nil {
		fmt.Printf("Failed to call ServerStream: %v\n", err)
		return
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Server stream completed")
			break
		}
		if err != nil {
			fmt.Printf("Received stream error: %v\n", err)
			break
		}

		fmt.Printf("Received from stream: %s\n", response.GetMessage())
	}
}

func callClientStream(client pb.SimpleServiceClient, scanner *bufio.Scanner) {
	rawMessages := readLine(scanner, "Enter messages separated by comma: ")
	errorParam := readBool(scanner, "Trigger error on last message? (y/n): ")

	parts := strings.Split(rawMessages, ",")
	var messages []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			messages = append(messages, trimmed)
		}
	}

	if len(messages) == 0 {
		messages = []string{"item-1", "item-2", "item-3"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ClientStream(ctx)
	if err != nil {
		fmt.Printf("Failed to open ClientStream: %v\n", err)
		return
	}

	for i, msg := range messages {
		isError := errorParam && (i == len(messages)-1)
		fmt.Printf("Sending incrementally: %q (error: %v)\n", msg, isError)

		err := stream.Send(&pb.Request{
			Message: msg,
			Error:   isError,
		})
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("Received close error: %v\n", err)
		return
	}

	fmt.Printf("Received final response: %s\n", response.GetMessage())
}

func callBiDirectionalStream(client pb.SimpleServiceClient, scanner *bufio.Scanner) {
	rawMessages := readLine(scanner, "Enter messages separated by comma: ")
	errorParam := readBool(scanner, "Trigger error on last message? (y/n): ")

	parts := strings.Split(rawMessages, ",")
	var messages []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			messages = append(messages, trimmed)
		}
	}

	if len(messages) == 0 {
		messages = []string{"msg-1", "msg-2", "msg-3"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.BiDirectionalStream(ctx)
	if err != nil {
		fmt.Printf("Failed to open BiDirectionalStream: %v\n", err)
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Bi-directional receive stream closed")
				return
			}
			if err != nil {
				fmt.Printf("Received stream error: %v\n", err)
				return
			}

			fmt.Printf("Received stream response: %s\n", response.GetMessage())
		}
	}()

	for i, msg := range messages {
		isError := errorParam && (i == len(messages)-1)
		fmt.Printf("Sending incrementally: %q (error: %v)\n", msg, isError)

		err := stream.Send(&pb.Request{
			Message: msg,
			Error:   isError,
		})
		if err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	if err := stream.CloseSend(); err != nil {
		fmt.Printf("Failed to close send stream: %v\n", err)
	}

	<-done
}

func printMenu() {
	fmt.Println("\n=== Select gRPC Request ===")
	fmt.Println("1. Normal")
	fmt.Println("2. Server Stream")
	fmt.Println("3. Client Stream")
	fmt.Println("4. Bi Directional")
	fmt.Println("5. Exit")
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleServiceClient(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()
		choice := readLine(scanner, "Select option (1-5): ")

		switch choice {
		case "1":
			callNormal(client, scanner)
		case "2":
			callServerStream(client, scanner)
		case "3":
			callClientStream(client, scanner)
		case "4":
			callBiDirectionalStream(client, scanner)
		case "5":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid option. Please enter 1, 2, 3, 4, or 5.")
		}
	}
}
