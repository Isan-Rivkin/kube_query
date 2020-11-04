BIN_DIR='/usr/local/bin'
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTOOL=$(GOCMD) tool
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD)fmt

GOARCH=$(shell go env GOARCH)

BINARY_NAME=kq
BINARY_LINUX=$(BINARY_NAME)_linux


# local development test
release:
	go build main.go
	mv main kq
	echo "Installing to ${BIN_DIR}" 
	mv kq ${BIN_DIR}
test:
	kq get pods -n kube-system
clean: 
	rm ${BIN_DIR}/kq

build-linux: ## Build Cross Platform Binary
		CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) $(GOBUILD) -o $(BINARY_NAME)_linux -v

build-osx: ## Build Mac Binary
		CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) $(GOBUILD) -o $(BINARY_NAME)_osx -v

build-windows: ## Build Windows Binary
		CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) $(GOBUILD) -o $(BINARY_NAME)_windows -v
