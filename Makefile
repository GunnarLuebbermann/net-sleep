# Makefile für NetSleep
APP_NAME = netsleep
VERSION  = 1.0.0

# Flags für kleinere Binary + Versionsinfo
LDFLAGS = -s -w -X 'main.version=$(VERSION)' -X 'main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'

.PHONY: all clean build

all: clean build

build:
	@echo "🚀 Baue Binaries..."
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)_windows_amd64.exe .
	GOOS=linux   GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)_linux_amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o build/$(APP_NAME)_mac_arm64 .
	@echo "✅ Fertig! Siehe ./build/"

clean:
	rm -rf build