		# Go parameters
    GOCMD=go
    GOFMT=$(GOCMD) fmt
    GOBINFLDR=
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=nerf-controller
    BINARY_UNIX=$(BINARY_NAME)_unix
    BINARY_EXTENTION=.exe

    # all: fmt lint test build
    all: fmt test build
    clean:
						$(GOCLEAN)
						rm -f $(BINARY_NAME)
						rm -f $(BINARY_UNIX)
    fmt:
						$(GOFMT)
    lint:
						$(GOGET) github.com/golang/lint/golint
						golint
    test:
						$(GOTEST) -v ./...
    build:
						$(GOBUILD) -o $(BINARY_NAME) -v $(GOBINFLDR)
    run:
						$(GOBUILD) -o $(BINARY_NAME) -v ./...
						./$(BINARY_NAME)
    # Cross compilation
    build-linux:
            CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
    docker-build:
            #docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
