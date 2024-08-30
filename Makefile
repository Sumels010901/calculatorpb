gen-cal:
	protoc --proto_path=./calculator/calculatorpb/ --go_out=./calculator/calculatorpb --go_opt=paths=source_relative --go-grpc_out=./calculator/calculatorpb/ --go-grpc_opt=paths=source_relative calculator/calculatorpb/calculator.proto
run-server:
	go run ./calculator/server/server.go
run-client:
	go run ./calculator/client/client.go