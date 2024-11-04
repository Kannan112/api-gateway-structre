# Variables
CMD_DIR=./cmd/server
PROTO_SCRIPT=scripts/generate-proto.sh

# Default target to run the Go application
run:
	go run $(CMD_DIR)/main.go

# Target to generate protobuf files
proto:
	bash $(PROTO_SCRIPT)

# You can add more targets as needed, e.g., build, clean, etc.
build:
	go build -o bin/server $(CMD_DIR)/main.go

# Clean the build output
clean:
	rm -rf bin/
