build:
	@go build -o bin/rofl

run: build
	@./bin/rofl

test:
	@go test -v ./...
