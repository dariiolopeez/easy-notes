BINARY_NAME := easy-notes
CMD         := ./cmd/easy-notes
BIN_DIR     := bin

.PHONY: build install test clean

build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD)

install:
	go install $(CMD)

test:
	go test -v ./...

clean:
	rm -rf $(BIN_DIR)
	go clean -testcache
