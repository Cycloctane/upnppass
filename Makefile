VERSION=$(shell git describe --tags --always)
LD_FLAGS=-s -w -X main.version=$(VERSION)
OUTPUT_DIR=build
WINDOWS=$(VERSION)_windows_x64.exe
LINUX=$(VERSION)_linux_x64

BINARY=upnpfind upnppass

.PHONY: all $(BINARY) clean

all: $(BINARY)

mkdir:
	@mkdir -p $(OUTPUT_DIR)

$(BINARY): mkdir
	GOOS=windows GOARCH=amd64 go build -v -o ./$(OUTPUT_DIR)/$@_$(WINDOWS) -ldflags="$(LD_FLAGS)" ./cmd/$@
	GOOS=linux GOARCH=amd64 go build -v -o ./$(OUTPUT_DIR)/$@_$(LINUX) -ldflags="$(LD_FLAGS)" ./cmd/$@

clean:
	rm $(OUTPUT_DIR)/*
