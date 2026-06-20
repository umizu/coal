build: 
	@go build -o bin/rail ./cmd
run: build
	@./bin/rail
