SHELL = /usr/bin/env bash

APP = go-vis
REPOSITORY = github.com/zchee/go-vis

GCFLAGS ?=
LDFLAGS ?=

PACKAGES = $(shell go list ./...)
FMT_PACKAGES = $(shell go list -f '{{.Dir}}' ./...)
GO_TEST_FLAGS ?=
GO_BENCH_FUNCS ?= .
GO_BENCH_FLAGS ?= $(GO_TEST_FLAGS) -run=^$$ -bench=${GO_BENCH_FUNCS}

all: $(APP)

$(APP): $(shell find . -name '*.go')
	go build -v $(GCFLAGS) $(LDFLAGS) -o ./bin/$(APP) $(REPOSITORY)/cmd/$(APP)

.PHONY: clean
clean:
	${RM} -r ./bin *.test *.out

init:
	go get -u -v github.com/golang/dep/cmd/dep

.PHONY: dep
dep:
	dep ensure -v


bin/goimports:
	go build -o bin/goimports ./vendor/golang.org/x/tools/cmd/goimports

bin/megacheck:
	go build -o bin/megacheck ./vendor/honnef.co/go/tools/cmd/megacheck

bin/errcheck:
	go build -o bin/errcheck ./vendor/github.com/kisielk/errcheck


.PHONY: test
test:
	go test -v -race $(GO_TEST_FLAGS) $(TEST_SRC_PACKAGES)

.PHONY: benchmark
benchmark-go: datastore-start
	go test -v $(GO_BENCH_FLAGS) $(TEST_SRC_PACKAGES)


.PHONY: lint
lint: fmt vet megacheck errcheck

.PHONY: fmt
fmt: bin/goimports
	gofmt -l $(FMT_PACKAGES) | grep -E '.'; test $$? -eq 1
	./bin/goimports -l $(FMT_PACKAGES) | grep -E '.'; test $$? -eq 1

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: staticcheck
megacheck: bin/megacheck
	./bin/megacheck $(PACKAGES)

.PHONY: errcheck
errcheck: bin/errcheck
	./bin/errcheck -exclude .errcheckignore $(PACKAGES)
