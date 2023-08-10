build:
	@go build -o bin/go-oci

run: build
	@./bin/go-oci

test:
	@go test -v ./...		