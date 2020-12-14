.PHONY: build clean 

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-gpio-moist-go
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-virtual-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/device-gpio-moist-go:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

clean:
	rm -f $(MICROSERVICES)
