BINARY_NAME=hotelapi
OSNAME=$(shell uname)

build:
ifeq ($(OSNAME),Linux)
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} main.go
else ifeq ($(OSNAME),Darwin)
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME} main.go
else
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME} main.go
endif

run: build
	./${BINARY_NAME}

clean:
	go clean
	@rm ${BINARY_NAME}

test:
	@go test -v ./...

test_coverage:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 api