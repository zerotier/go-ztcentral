GOLANGCI_LINT_VERSION := 1.34.1

lint: bin/golangci-lint
	bin/golangci-lint run -v

reflex-lint: bin/reflex
	bin/reflex -r '\.go$$' make lint

reflex-test: bin/reflex
	bin/reflex -r '\.go$$' -- go test -v ./...

bin/golangci-lint:
	mkdir -p bin
	wget -O- https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz | tar vxz --strip-components=1 -C bin golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64/golangci-lint

bin/reflex:
	GO111MODULE=off GOBIN=${PWD}/bin go get -u github.com/cespare/reflex
