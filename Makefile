BINARY_NAME=chess2pic

build: chess2pic

chess2pic:
	go build -o build/${BINARY_NAME} ./app/chess2pic

clean:
	go clean

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...