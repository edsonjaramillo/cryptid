cli:
	go run backend/cmd/cli.go

lint:
	golangci-lint run 

format:
	golangci-lint fmt
