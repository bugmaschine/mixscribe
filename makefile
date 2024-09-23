BUILD_DIR = bin

all: server

server:
	@echo "Building Server"
	GOOS=linux GOARCH=amd64 go build -a -tags osusergo,netgo -ldflags '-w -extldflags "-static"' -o $(BUILD_DIR)/server .
	GOOS=windows GOARCH=amd64 go build -a -tags osusergo,netgo -ldflags '-w -extldflags "-static"' -o $(BUILD_DIR)/server.exe .
