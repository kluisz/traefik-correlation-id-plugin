.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	$(eval GOPATH_TMP := $(shell mktemp -d))
	mkdir -p $(GOPATH_TMP)/src/github.com/kluisz
	ln -s $(CURDIR) $(GOPATH_TMP)/src/github.com/kluisz/traefik-correlation-id-plugin
	GOPATH=$(GOPATH_TMP) yaegi test -v .
	rm -rf $(GOPATH_TMP)

vendor:
	go mod vendor

clean:
	rm -rf ./vendor