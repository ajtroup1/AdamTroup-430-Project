# Define the Go binary name and source directory
APP_NAME = godoc
SRC_DIR = .

# Define the default task
all: build

# Compile the Go program
build:
	go build -o $(APP_NAME) $(SRC_DIR)/bin

# Run the program with the 'save' task
save:
	go run cmd/main.go -task save

# Run the program with the 'gen' task
gen:
	go run cmd/main.go -task gen

# Clean the binary
clean:
	rm -f $(APP_NAME)

fmt:
	go fmt ./...