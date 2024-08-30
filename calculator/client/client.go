package main

import (
	"calculator/calculator/calculatorpb"
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("Client calling sum api...")
	response, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 254,
		Num2: 19,
	})

	if err != nil {
		log.Fatalf("Error calling sum: %v", err)
	}

	log.Printf("Sum api response: %v\n", response)
}

func callSumWithDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	log.Println("Client calling sum with deadline api...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		Num1: 254,
		Num2: 19,
	})

	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			if errStatus.Code() == codes.DeadlineExceeded {
				log.Println("Deadline exceeded, cancel the request...")
			} else {
				log.Printf("Error calling sum with deadline: %v", err)
			}
		} else {
			log.Printf("Error calling sum with deadline: %v", err)
		}
		log.Printf("Error unknown while calling sum with deadline: %v", err)
		return
	}

	log.Printf("Sum api response: %v\n", response.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("Client calling pnd api...")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("Call PND err: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Server finished streaming: %v", err)
			return
		}

		log.Printf("Prime number: %v", response.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("Client calling avg api...")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Average error: %v", err)
	}

	listRequest := []calculatorpb.AverageRequest{
		{
			Num: 5,
		},
		{
			Num: 10,
		},
		{
			Num: 15,
		},
		{
			Num: 20,
		},
		{
			Num: 25,
		},
		{
			Num: 30,
		},
	}

	for _, request := range listRequest {
		err := stream.Send(&request)
		if err != nil {
			log.Printf("Error sending request %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error receiving response %v", err)
	}

	log.Printf("Average response received %v", response)
}

func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("Client calling find max api...")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Average error: %v", err)
	}

	listRequest := []calculatorpb.MaxRequest{
		{
			Num: 625415,
		},
		{
			Num: 151230,
		},
		{
			Num: 412341,
		},
		{
			Num: 123123,
		},
		{
			Num: 213442,
		},
		{
			Num: 382748,
		},
	}

	waitc := make(chan struct{})

	go func() {
		// Handle multiple requests
		for _, request := range listRequest {
			err := stream.Send(&request)
			if err != nil {
				log.Printf("Error sending request %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				log.Println("Ending find max api...")
				break
			}
			if err != nil {
				log.Printf("Error receiving response %v", err)
				break
			}
			log.Printf("Max : %v", response.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

func callSqrt(c calculatorpb.CalculatorServiceClient, num int32) {
	log.Println("Client calling sqrt api...")
	response, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})

	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("Error message: %v\n", errStatus.Message())
			log.Printf("Error code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("Invalid number < 0: %v", num)
				return
			}
		}
	}

	log.Printf("Sqrt api response: %v\n", response.GetSqrt())
}

func main() {
	conn, err := grpc.Dial("localhost:19110", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while dial: %v", err)
	}

	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	//callSum(client)
	callSumWithDeadline(client, 1*time.Second)
	callSumWithDeadline(client, 5*time.Second)
	//callPND(client)
	//callAverage(client)
	//callFindMax(client)
	//callSqrt(client, -9)
}
