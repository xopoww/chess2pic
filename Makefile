BINARY_NAME=chess2pic

build: chess2pic chess2pic-api

chess2pic:
	go build -o build/${BINARY_NAME} ./cmd/chess2pic

chess2pic-api:
	swagger generate server -A chess2pic-api -f api/chess2pic-api.yml
	go build -o build/${BINARY_NAME} ./cmd/chess2pic-api-server

clean:
	go clean

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...