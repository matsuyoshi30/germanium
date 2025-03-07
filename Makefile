GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=germanium

all: build test
build:
	$(GOBUILD) -o dist/$(BINARY_NAME) -v ./cmd/$(BINARY_NAME)
build_windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/$(BINARY_NAME).exe -v ./cmd/$(BINARY_NAME)	
test:
	$(GOTEST) -v ./...
gentest: build
	$(GOTEST) -v ./... -gen_golden_files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
