#!/usr/bin/make -f
SRC_DIR	:= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

AUTOGRAF_BIN := ./autograf

all: deps lint build goreleaser-build test e2e-test

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

RELEASE_NOTES ?= $(TMP_DIR)/release_notes
$(RELEASE_NOTES): $(TMP_DIR)
	@echo "Generating release notes to $(RELEASE_NOTES) ..."
	@csplit -q -n1 --suppress-matched -f $(TMP_DIR)/release-notes-part CHANGELOG.md '/## \[\s*v.*\]/' {1}
	@mv $(TMP_DIR)/release-notes-part1 $(RELEASE_NOTES)
	@rm $(TMP_DIR)/release-notes-part*

lint:
	golangci-lint run

test:
	go test -race ./...

e2e-test: build
	$(AUTOGRAF_BIN) --metrics-file examples/metrics.txt --grafana-url=""

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(AUTOGRAF_BIN)

goreleaser-build:
	goreleaser build --snapshot --rm-dist

docker: build
	docker build -t fusakla/autograf .

.PHONY: clean
clean:
	rm -rf dist $(AUTOGRAF_BIN) $(TMP_DIR)

.PHONY: deps
deps:
	go mod tidy && go mod verify
