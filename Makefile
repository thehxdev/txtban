GO := go
BIN := tb

$(BIN): $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=1 $(GO) build -ldflags='-s' -o $(BIN) .

clean:
	rm -rf $(BIN)
