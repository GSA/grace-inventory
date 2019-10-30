GOOS = linux
GOARCH = amd64
GOBIN := $(GOPATH)/bin
GODEP := $(GOBIN)/dep
GOLANGCILINT := $(GOBIN)/golangci-lint
GOSEC := $(GOBIN)/gosec
GHR := $(GOBIN)/ghr
RELEASEDIR ?= ../release/
CIRCLE_PROJECT_REPONAME ?= grace-inventory-lambda

.PHONY: build release clean test lint dependencies
build: lint clean
	mkdir -p $(RELEASEDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME) -v
	zip -j $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME).zip $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME)
	rm -f $(RELEASEDIR)$(CIRCLE_PROJECT_REPONAME)

release: build
	export CIRCLE_TAG=$(CIRCLE_TAG)
ifeq ($(strip $(GITHUB_TOKEN)),)
	@echo "GITHUB_TOKEN must be provided"
	@exit 0
else ifeq ($(strip $(CIRCLE_TAG)),)
  @echo "CIRCLE_TAG must be set"
	@exit 0
else ifeq ($(shell echo $(CIRCLE_TAG) | egrep '^v\d+\.\d+\.\d+'),)
  @echo "CIRCLE_TAG must match version pattern (i.e. v.1.2.3)"
	@exit 0
else
	ghr -t $(GITHUB_TOKEN) -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) -c $(CIRCLE_SHA1) -delete $(CIRCLE_TAG) $(RELEASEDIR)
endif

clean:
	rm -rf $(RELEASEDIR)

test: lint
	go test -v ./...

test_nocreds: lint
	go test -v --run 'Users|Groups|Policies|Roles|Vaults|TestInstances|Images|Volumes|TestSnapshots|Vpcs|Subnets|Addresses|KeyPairs|Buckets' ./...

lint: Gopkg.toml dependencies
	dep ensure
	golangci-lint run ./...
	gosec ./...

Gopkg.toml:
ifeq (,$(wildcard Gopkg.toml))
	dep init
endif

dependencies: $(GODEP) $(GOLANGCILINT) $(GOSEC) $(GHR)

$(GOLANGCILINT):
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

$(GODEP):
	go get -u github.com/golang/dep/cmd/dep

$(GOSEC):
	go get -u github.com/securego/gosec/cmd/gosec

$(GHR):
	go get -u github.com/tcnksm/ghr
