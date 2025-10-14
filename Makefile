MAIN_NAME=dcvix-stats
BINARY_NAME=dcvix-Stats
GO_PACKAGE=github.com/dcvix/$(MAIN_NAME)

# Version information
VERSION=$(shell cat VERSION)
RELEASE=1
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build variables
GO=go
LDFLAGS="-X $(GO_PACKAGE)/internal/version.Version=$(VERSION) \
         -X $(GO_PACKAGE)/internal/version.Commit=$(COMMIT) \
         -X $(GO_PACKAGE)/internal/version.BuildTime=$(BUILD_TIME)"
LDFLAGS_WIN=-X=$(GO_PACKAGE)/internal/version.Version=$(VERSION),-X=$(GO_PACKAGE)/internal/version.Commit=$(COMMIT),-X=$(GO_PACKAGE)/internal/version.BuildTime=$(BUILD_TIME)

# Platform-specific variables
LINUX_AMD64_BINARY=$(MAIN_NAME)
LINUX_AMD64_DIR=$(MAIN_NAME)-v$(VERSION)-linux-amd64
WINDOWS_AMD64_BINARY=$(MAIN_NAME).exe
WINDOWS_AMD64_DIR=$(MAIN_NAME)-v$(VERSION)-windows-amd64

# Build all platforms
.PHONY: build
build: build-linux build-windows-cross

# Build for Linux
.PHONY: build-linux
build-linux: update-toml
	mkdir -p dist/$(LINUX_AMD64_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o dist/$(LINUX_AMD64_DIR)/$(LINUX_AMD64_BINARY) ./cmd/$(MAIN_NAME)
	cp README.md LICENSE.md dist/$(LINUX_AMD64_DIR)/
	cd dist && tar czf $(LINUX_AMD64_DIR).tar.gz $(LINUX_AMD64_DIR)

# Build for Windows
.PHONY: build-windows
build-windows: update-toml
	go-winres simply --product-version $(VERSION).0 --file-version $(VERSION).0 --file-description "Graphical interface to easily launch the DCV viewer" --product-name "DCV Launcher" --copyright "Diego Cortassa" --original-filename "$(WINDOWS_BINARY)" --icon Icon.png
	mkdir -p dist/$(WINDOWS_AMD64_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o dist/$(WINDOWS_AMD64_DIR)/$(WINDOWS_BINARY) ./cmd/$(MAIN_NAME)
	cp README.md LICENSE.md dist/$(WINDOWS_AMD64_DIR)/
	cd dist && 7z a -r $(WINDOWS_AMD64_DIR).zip $(WINDOWS_AMD64_DIR)

# Build for Windows cross compile
.PHONY: build-windows-cross
build-windows-cross: update-toml
	go-winres simply --product-version $(VERSION).0 --file-version $(VERSION).0 --file-description "Graphical interface to easily launch the DCV viewer" --product-name "DCV Launcher" --copyright "Diego Cortassa" --original-filename "$(WINDOWS_BINARY)" --icon Icon.png
	GOFLAGS="-ldflags=$(LDFLAGS_WIN)" fyne-cross windows -arch=amd64 -icon=Icon.png ./cmd/$(MAIN_NAME)
	mkdir -p dist/$(WINDOWS_AMD64_DIR)	
	mv fyne-cross/bin/windows-amd64/$(MAIN_NAME).exe dist/$(WINDOWS_AMD64_DIR)/$(WINDOWS_AMD64_BINARY)
	cp README.md LICENSE.md dist/$(WINDOWS_AMD64_DIR)/
	cd dist && 7z a -r $(WINDOWS_AMD64_DIR).zip $(WINDOWS_AMD64_DIR)

# Generate embeds
PHONY: generate
generate:
	go generate internal/gui/gui.go ;

# Print version
.PHONY: version
version:
	@echo $(VERSION)

# Bump version (patch by default)
.PHONY:
version-bump: 
	@current_version=`cat VERSION`; \
	major=`echo $$current_version | cut -d. -f1`; \
	minor=`echo $$current_version | cut -d. -f2`; \
	patch=`echo $$current_version | cut -d. -f3`; \
	new_minor=$$((minor + 1)); \
	new_version="$$major.$$new_minor.$$patch"; \
	echo $$new_version > VERSION; \
	echo "Version bumped from $$current_version to $$new_version"
	$(MAKE) update-toml

# update-toml
.PHONY: update-toml
update-toml:
	sed -i "s/^  Version = \".*\"/  Version = \"$(VERSION)\"/" FyneApp.toml
	sed -i "s/^  Build = .*/  Build = $(RELEASE)/" FyneApp.toml

# Create a new version tag
.PHONY: tag
tag: version
	git add VERSION FyneApp.toml
	@if git diff --quiet --cached -- VERSION FyneApp.toml; then \
		echo "VERSION and FyneApp.toml up to date, tagging"; \
		git tag -a v$(VERSION) -m "Version $(VERSION)"; \
		echo "Tagged, now push to GitHub: git push origin v$(VERSION)"; \
	else \
		echo "VERSION and FyneApp.toml need to be committed first"; \
	fi

PHONY: clean
clean:
	# go clean ;
	rm -rf dist
	rm -rf fyne-cross
	rm -f *.syso

PHONY: run-debug
run-debug:
	go run -tags debug cmd/$(MAIN_NAME)/main.go --verbose --logfile examples/server.log ;

PHONY: run
run:
	go run cmd/$(MAIN_NAME)/main.go --verbose --logfile examples/server.log ;
