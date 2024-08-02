GOLANGCI_LINT_VERSION := 1.59.1
VERSIONFILE=VERSION

lint: docker-pull
	docker run --rm -v $(CURDIR):/app -v /tmp/.cache/golangci-lint/v1.59.1:/root/.cache -w /app golangci/golangci-lint:v1.59.1 golangci-lint run -v

docker-pull:
	docker pull golangci/golangci-lint:v${GOLANGCI_LINT_VERSION}

reflex-lint: bin/reflex
	bin/reflex -r '\.go$$' make lint

reflex-test: bin/reflex
	bin/reflex -r '\.go$$' -- go test -v ./...

bin/reflex:
	GOBIN=${PWD}/bin go get -u github.com/cespare/reflex

generate:
	go generate -v ./...

update-openapi:
	curl -sSL https://raw.githubusercontent.com/zerotier/docs/master/static/openapi/centralv1.json >spec.json
	go generate ./...

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
