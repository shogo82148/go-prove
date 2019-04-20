GOVERSION=$(shell go version)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const Version' version.go)))
ARTIFACTS_DIR=$(CURDIR)/artifacts/$(VERSION)
RELEASE_DIR=$(CURDIR)/release/$(VERSION)
SRC_FILES = $(wildcard *.go cli/go-prove/*.go formatter/*.go plugin/*.go)
GITHUB_USERNAME=shogo82148
ARCHIVER=$(CURDIR)/archiver-$(shell go env GOOS)-$(shell go env GOARCH)/archiver

.PHONY: all test clean

all: build-windows-386 build-windows-amd64 build-linux-386 build-linux-amd64 build-darwin-386 build-darwin-amd64

##### build settings

.PHONY: build build-windows-amd64 build-windows-386 build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386

$(ARTIFACTS_DIR)/go-prove_$(GOOS)_$(GOARCH):
	@mkdir -p $@

$(ARTIFACTS_DIR)/go-prove_$(GOOS)_$(GOARCH)/go-prove$(SUFFIX): $(ARTIFACTS_DIR)/go-prove_$(GOOS)_$(GOARCH) $(SRC_FILES)
	@echo " * Building binary for $(GOOS)/$(GOARCH)..."
	@CGO_ENABLED=0 go build -o $@ cli/go-prove/main.go

build: $(ARTIFACTS_DIR)/go-prove_$(GOOS)_$(GOARCH)/go-prove$(SUFFIX)

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386 SUFFIX=.exe

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

##### release settings

.PHONY: release-windows-amd64 release-windows-386 release-linux-amd64 release-linux-386 release-darwin-amd64 release-darwin-386
.PHONY: release-targz release-zip release-files release-upload

$(RELEASE_DIR)/go-prove_$(GOOS)_$(GOARCH):
	@mkdir -p $@

release-windows-amd64:
	@$(MAKE) release-zip GOOS=windows GOARCH=amd64 SUFFIX=.exe

release-windows-386:
	@$(MAKE) release-zip GOOS=windows GOARCH=386 SUFFIX=.exe

release-linux-amd64:
	@$(MAKE) release-targz GOOS=linux GOARCH=amd64

release-linux-386:
	@$(MAKE) release-targz GOOS=linux GOARCH=386

release-darwin-amd64:
	@$(MAKE) release-zip GOOS=darwin GOARCH=amd64

release-darwin-386:
	@$(MAKE) release-zip GOOS=darwin GOARCH=386

release-targz: build $(RELEASE_DIR)/go-prove_$(GOOS)_$(GOARCH)
	@echo " * Creating tar.gz for $(GOOS)/$(GOARCH)"
	tar -czf $(RELEASE_DIR)/go-prove_$(GOOS)_$(GOARCH).tar.gz -C $(ARTIFACTS_DIR) go-prove_$(GOOS)_$(GOARCH)

release-zip: build $(RELEASE_DIR)/go-prove_$(GOOS)_$(GOARCH)
	@echo " * Creating zip for $(GOOS)/$(GOARCH)"
	cd $(ARTIFACTS_DIR) && zip -9 $(RELEASE_DIR)/go-prove_$(GOOS)_$(GOARCH).zip go-prove_$(GOOS)_$(GOARCH)/*

release-files: release-windows-386 release-windows-amd64 release-linux-386 release-linux-amd64 release-darwin-386 release-darwin-amd64

release-upload: release-files
	ghr -u $(GITHUB_USERNAME) --draft --replace v$(VERSION) $(RELEASE_DIR)

$(ARCHIVER):
	mkdir -p $(shell dirname $(ARCHIVER))
	wget -O $(ARCHIVER) https://github.com/mholt/archiver/releases/download/v2.0/archiver_$(shell go env GOOS)_$(shell go env GOARCH)
	chmod 755 $(ARCHIVER)

test:
	go test -v -race ./...
	go vet ./...

clean:
	-rm -rf vendor
