build: chess2pic chess2pic-api

chess2pic:
	go build -o build/chess2pic ./cmd/chess2pic

chess2pic-api:
	go generate ./restapi
	go mod tidy -compat=1.17
	go build -o build/chess2pic-api-server ./cmd/chess2pic-api-server

run-api:
	PORT=65000 ./build/chess2pic-api-server

clean:
	go clean

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...