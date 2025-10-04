MAIN_NAME=dcvix-stats

# Version information
VERSION=$(shell cat VERSION)
RELEASE=1
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build variables
BINARY_NAME=$(MAIN_NAME)
GO=go
GOFMT=gofmt
GOFILES=$(shell find . -name "*.go")
LDFLAGS="-X github.com/dcvix/$(MAIN_NAME)/internal/version.Version=$(VERSION) \
         -X github.com/dcvix/$(MAIN_NAME)/internal/version.Commit=$(COMMIT) \
         -X github.com/dcvix/$(MAIN_NAME)/internal/version.BuildTime=$(BUILD_TIME)"
LDFLAGS_WIN=-X=github.com/dcvix/$(MAIN_NAME)/internal/version.Version=$(VERSION),-X=github.com/dcvix/$(MAIN_NAME)/internal/version.Commit=$(COMMIT),-X=github.com/dcvix/$(MAIN_NAME)/internal/version.BuildTime=$(BUILD_TIME)

# Platform-specific variables
WINDOWS_BINARY=$(BINARY_NAME).exe
WINDOWS_AMD64_DIR=dist
LINUX_AMD64_DIR=$(WINDOWS_AMD64_DIR)

# Build all platforms
.PHONY: build
build: update-toml build-linux build-windows-cross

# update-toml
.PHONY: update-toml
update-toml:
	sed -i "s/^  Version = \".*\"/  Version = \"$(VERSION)\"/" FyneApp.toml
	sed -i "s/^  Build = .*/  Build = $(RELEASE)/" FyneApp.toml

# Build for Linux
.PHONY: build-linux
build-linux:
	mkdir -p $(LINUX_AMD64_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(LINUX_AMD64_DIR)/$(BINARY_NAME) ./cmd/$(MAIN_NAME)

# Build for Windows
.PHONY: build-windows
build-windows:
	mkdir -p $(WINDOWS_AMD64_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o $(WINDOWS_AMD64_DIR)/$(WINDOWS_BINARY) ./cmd/$(MAIN_NAME)

# Build for Windows cross compile
.PHONY: build-windows-cross
build-windows-cross:
	GOFLAGS="-ldflags=$(LDFLAGS_WIN)" fyne-cross windows -arch=amd64 -icon=./assets/icon.png ./cmd/dcvix-stats
	mv fyne-cross/bin/windows-amd64/dcvix-stats.exe $(LINUX_AMD64_DIR)/$(WINDOWS_BINARY)

# Show version
.PHONY: version
version:
	@echo $(VERSION)

# Create a new version tag
.PHONY: tag
tag: version update-toml
	git tag -a v$(VERSION) -m "Version $(VERSION)"
	# git push origin v$(VERSION)

PHONY: clean
clean:
	# go clean ;
	rm -rf dist
	rm -rf fyne-cross

PHONY: run
run:
# 	go run -tags debug cmd/dcvix-stats/main.go --verbose --logfile examples/server.log ;
# 	go run cmd/dcvix-stats/main.go --verbose --logfile examples/server.log ;
	go run cmd/dcvix-stats/main.go --logfile examples/server.log ;
