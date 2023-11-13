ldflags=-X=github.com/ipfs-force-community/api-compare/version.CurrentCommit=+git.$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))
ifneq ($(strip $(LDFLAGS)),)
	    ldflags+=-extldflags=$(LDFLAGS)
	endif

GOFLAGS+=-ldflags="$(ldflags)"

build: api-compare

api-compare:
	rm -rf api-compare
	go build -o api-compare $(GOFLAGS) main.go
.PHONY: api-compare

lint:
	golangci-lint run

test:
	go test -v ./...
