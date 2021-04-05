GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=germanium

all: build test
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test: build
	$(GOTEST) -v ./...
gentest:
	$(GOTEST) -gen_golden_files ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
