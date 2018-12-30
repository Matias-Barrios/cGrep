# Go parameters
GOCMD=go 
GOBUILD=$(GOCMD) build
BINARY_NAME=cgrep
BINARY_NAME_MACOS=cgrep_darwin

clean: 
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME_MACOS)
   
# Specific compilation

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME_MACOS) -v