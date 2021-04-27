GOLANGCI_LINT_VERSION := 1.34.1
VERSIONFILE=VERSION

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

generate:
	go generate -v ./...

update-openapi:
	curl -sSL http://apidocs.zerotier.com/central-v1/api-spec.json >spec.json

release:
ifeq (${VERSION},)
	@echo "Please set VERSION before running this make task"
	exit 1
endif
	@# our generate process requires a tagged version, otherwise it will append a
	@# sha; this retagging mess works around that. Computers.
	git tag "${VERSION}" # git tag -d if this fails
	@echo "${VERSION}" > "$(VERSIONFILE)"
	@make generate
	@git tag -d "${VERSION}" # strip the tag so we can paper on our commit
	git commit -a -s -m "Release ${VERSION}"
	git tag "${VERSION}" # reincorporate the tag
	@echo "Please push tag '${VERSION}' to the appropriate remote."
