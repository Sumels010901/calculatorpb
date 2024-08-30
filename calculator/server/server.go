package main

import (
	"calculator/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

func (*server) Sum(ctx context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum called...")
	response := &calculatorpb.SumResponse{
		Result: request.GetNum1() + request.GetNum2(),
	}

	return response, nil
}

func (*server) SumWithDeadline(ctx context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum with deadline called...")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("Client canceled request")
			return nil, status.Errorf(codes.Canceled, "client canceled request")
		}
		time.Sleep(1 * time.Second)
	}

	response := &calculatorpb.SumResponse{
		Result: request.GetNum1() + request.GetNum2(),
	}

	return response, nil
}

func (*server) PrimeNumberDecomposition(request *calculatorpb.PNDRequest,
	stream grpc.ServerStreamingServer[calculatorpb.PNDResponse]) error {
	log.Println("Prime Number Decomposition called...")
	k := int32(2)
	N := int32(request.GetNumber())
	for N > 1 {
		if N%k == 0 {
			N = N / k
			// send to client
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %d", k)
		}
	}
	return nil
}

func (*server) Average(stream grpc.ClientStreamingServer[calculatorpb.AverageRequest, calculatorpb.AverageResponse]) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Server finished streaming: %v", err)
			response := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(response)
		}
		if err != nil {
			log.Fatalf("Error streaming: %v", err)
			return err
		}
		log.Printf("Received request %v", request)
		total += request.GetNum()
		count++
	}
}

func (*server) FindMax(stream grpc.BidiStreamingServer[calculatorpb.MaxRequest, calculatorpb.MaxResponse]) error {
	log.Println("FindMax called...")
	max := int32(0)
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF reached")
			return nil
		}
		if err != nil {
			log.Fatalf("Error streaming: %v", err)
			return err
		}

		num := request.GetNum()
		log.Printf("Number received from client: %d", num)
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.MaxResponse{
			Max: max,
		})

		if err != nil {
			log.Fatalf("Error sending: %v", err)
			return err
		}
	}
}

func (*server) SquareRoot(ctx context.Context, request *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Println("SquareRoot called...")
	num := request.GetNum()
	if num < 0 {
		log.Printf("Invalid number %v less then 0", num)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid number %v less then 0", num)
	}
	return &calculatorpb.SquareResponse{
		Sqrt: math.Sqrt(float64(num)),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:19110")
	if err != nil {
		log.Fatalf("Error while create listener: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Calculator service is running...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Error while serve: %v", err)
	}
}
