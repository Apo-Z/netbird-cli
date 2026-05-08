.PHONY: build install clean

BINARY = netbird-cli
CMD    = ./cmd/cli

build:
	go build -ldflags="-s -w" -o $(BINARY) $(CMD)

install: build
	sudo install -m 755 $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)
