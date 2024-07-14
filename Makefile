GO := go
BIN := tb
DATA_DIR := ./data

$(BIN): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=1 $(GO) build -ldflags='-s' -o $(BIN) .

# Build docker image
docker: clean ./Dockerfile
	docker build -t txtban .

# Build executable file in golang docker container
# Using bullseye version because of glibc backward compatibility
docker_exe:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.22-bullseye make

clean:
	rm -rf *.db $(BIN) $(DATA_DIR)
