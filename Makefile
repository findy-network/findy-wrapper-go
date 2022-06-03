COV_FILE:=coverage.txt

deps:
	go get -t ./...

build:
	go build ./...

vet:
	go vet ./...

shadow:
	@echo Running govet
	go vet -vettool=$(GOPATH)/bin/shadow ./...
#	$(GOPATH)/bin/shadow ./...
	@echo Govet success

check_fmt:
	$(eval GOFILES = $(shell find . -name '*.go'))
	@gofmt -l $(GOFILES)

lint_old:
	$(GOPATH)/bin/golint ./...

lint_e:
	@$(GOPATH)/bin/golint ./... | grep -v export | cat

testv:
	go test -v -p 1 -failfast ./...

test:
	go test -p 1 -failfast ./...

ledger_test:
	go test -v -p 1 -failfast ./anoncreds/... -args -logtostderr=true -v=10

logged_test:
	go test -v -p 1 -failfast ./... -args -logtostderr=true -v=10

test_cov_out:
	go test \
		-coverpkg=github.com/findy-network/findy-wrapper-go/... \
		-coverprofile=$(COV_FILE)  \
		-covermode=atomic \
		./...

test_cov: test_cov_out
	go tool cover -html=$(COV_FILE)

check: check_fmt vet shadow

lint: lint_ci

lint_ci:
	golangci-lint run ./...

indy_to_debian:
	./scripts/debian-libindy/install-indy.sh
