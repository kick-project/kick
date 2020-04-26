SHELL = /bin/bash

# Project
NAME := $(shell basename $$(pwd))
GOPATH := $(shell go env GOPATH)
VERSION ?= $(shell cat VERSION)
COMMIT := $(shell test -d .git && git rev-parse --short HEAD)
BUILD_INFO := $(COMMIT)-$(shell date -u +"%Y%m%d-%H%M%SZ")
HASCMD := $(shell test -d cmd && echo "true")
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | perl -p -e 's/x86_64/amd64/; s/i386/386/')

RELEASE_DESCRIPTION := Standard release

ISRELEASED := $(shell git show-ref v$$(cat VERSION) 2>&1 > /dev/null && echo "true")

# Utilities
# Default environment variables.
# Any variables already set will override the values in this file(s).
DOTENV := godotenv -f $(HOME)/.env,.env

# Variables
ROOT = $(shell pwd)

# Go
GOMODOPTS = GO111MODULE=on
GOGETOPTS = GO111MODULE=off
GOFILES := $(shell find cmd pkg internal src -name '*.go' 2> /dev/null)
GODIRS = $(shell find . -maxdepth 1 -mindepth 1 -type d | egrep 'cmd|internal|pkg|api')

.PHONY: build _build _build_xcompile browsetest cattest clean deps _deps depsdev deploy _go.mod _go.mod_err \
        _isreleased lint release _release _release_gitlab test _test _test_setup _test_setup_home _test_setup_gitserver \
		unit codecomplexity codecoverage _unit _codecomplexity _codecoverage 

#
# End user targets
#
build: deps
	@VERSION=$(VERSION) $(DOTENV) make _build

buildauto:
	fswatch -o cmd/* internal/* --one-per-batch | xargs -n1 -I{} bash -c 'echo "make build # $$(date)"; make build; echo'

install:
	@VERSION=$(VERSION) $(DOTENV) make _install

installauto:
	fswatch -o cmd/* internal/* --one-per-batch | xargs -n1 -I{} bash -c 'echo "make install # $$(date)"; make install; echo'

clean:
	-test -f tmp/server.pid && kill -TERM $$(cat tmp/server.pid)
	rm -rf .cache prjstart dist reports tmp vendor

testsetup:
	@$(DOTENV) make _test_setup

test: deps
	@$(DOTENV) make _unit
	@$(DOTENV) make _codecoverage
	@$(DOTENV) make _codecomplexity

unit: deps
	@$(DOTENV) make _unit

codecoverage: 
	@$(DOTENV) make _codecoverage

codecomplexity:
	@$(DOTENV) make _codecomplexity

testauto:
	fswatch -o cmd/* internal/* test/* --one-per-batch | xargs -n1 -I{} bash -c 'echo "make test # $$(date)"; make test; echo'

lint:
	golangci-lint run --enable=gocyclo

deploy: build
	@echo TODO

release:
	@VERSION=$(VERSION) $(DOTENV) make _release 2> /dev/null

release_publish:
	@VERSION=$(VERSION) $(DOTENV) goreleaser release

.PHONY: tag
tag:
	git fetch --tags
	git tag v$(VERSION)
	git push --tags

deps: go.mod
ifeq ($(USEGITLAB),true)
	@mkdir -p $(ROOT)/.cache/{go,gomod}
endif
	@$(DOTENV) make _deps

depsdev:
ifeq ($(USEGITLAB),true)
	@mkdir -p $(ROOT)/.cache/{go,gomod}
endif
	@make $(GOGETS)

bumpmajor:
	git fetch --tags
	versionbump --checktags major VERSION

bumpminor:
	git fetch --tags
	versionbump --checktags minor VERSION

bumppatch:
	git fetch --tags
	versionbump --checktags patch VERSION

browsetest:
	@make $(REPORTS)

cattest:
	### Unit Tests
	@cat reports/test.txt
	### Code Coverage
	@cat reports/coverage.txt
	### Cyclomatix Complexity Report
	@cat reports/cyclocomplexity.txt

.PHONY: getversion
getversion:
	VERSION=$(VERSION) bash -c 'echo $$VERSION'

#
# Helper targets
#
_build:
	@test -d .cache || go fmt ./...
ifeq ($(HASCMD),true)
	@make $(NAME)
endif

.PHONY: _install
_install: $(GOPATH)/bin/$(NAME)

$(GOPATH)/bin/$(NAME): $(NAME)
	install -m 755 $(NAME) $(GOPATH)/bin/$(NAME)

_deps:
	$(GOMODOPTS) go mod tidy
	$(GOMODOPTS) go mod vendor

GOGETS := github.com/jstemmer/go-junit-report github.com/golangci/golangci-lint/cmd/golangci-lint \
		  github.com/ains/go-test-html github.com/goreleaser/goreleaser github.com/fzipp/gocyclo github.com/joho/godotenv/cmd/godotenv \
		  github.com/crosseyed/versionbump/cmd/versionbump github.com/stretchr/testify github.com/sosedoff/gitkit
.PHONY: $(GOGETS)
$(GOGETS):
	go get -u $@

_unit: _test_setup
	@make _test_setup_gitserver
	### Unit Tests
	@(go test -timeout 5s -covermode atomic -coverprofile=./reports/coverage.out -v ./...; echo $$? > reports/exitcode.txt) 2>&1 | tee reports/test.txt
	@cat ./reports/test.txt | go-junit-report > reports/junit.xml
	@exit $$(cat reports/exitcode.txt)

_codecoverage: _test_setup
	### Code Coverage
	@go tool cover -func=./reports/coverage.out | tee ./reports/coverage.txt
	@go tool cover -html=reports/coverage.out -o reports/html/coverage.html

_codecomplexity: _test_setup
	### Cyclomatix Complexity Report
	@gocyclo -avg $(GODIRS) | grep -v _test.go | tee reports/cyclocomplexity.txt

_test_setup:
	@mkdir -p tmp
	@mkdir -p reports/html
	@make _test_setup_home 2> /dev/null > /dev/null
	@sync

_test_setup_home:
	@cp -r test/fixtures/home tmp/

_test_setup_gitserver:
	@mkdir -p tmp/gitserveclient
	@-kill -0 $$(cat tmp/server.pid) 2>/dev/null >/dev/null || go run test/fixtures/testserver.go
	@echo "Waiting for git server to launch on 5000..."
	@bash -c 'while ! nc -z localhost 5000; do sleep 0.1; done'
	@echo "git server launched"
	@-find test/fixtures/gitserve -mindepth 1 -maxdepth 1 -type d -exec cp -r {} tmp/gitserveclient \;
	@-for i in $$(pwd)/tmp/gitserveclient/*; do cd $$i; git init; git add .; git commit -m "Initial commit"; git push --set-upstream http://127.0.0.1:5000/$$(basename $$(pwd)).git master; done
	@sync

_release:
	@echo "### Releasing v$(VERSION)"
	@make --no-print-directory _isreleased 2> /dev/null
	git tag v$(VERSION)
	git push --tags

REPORTS = reports/html/coverage.html
.PHONY: $(REPORTS)
$(REPORTS):
	@test -f $@ && open $@

# Check versionbump
_isreleased:
ifeq ($(ISRELEASED),true)
	@echo "Version $(VERSION) has been released."
	@echo "Please bump with 'make bump(minor|patch|major)' depending on breaking changes."
	@exit 1
endif

#
# File targets
#
$(NAME): dist/$(NAME)_$(OS)_$(ARCH)/$(NAME)
	install -m 755 dist/$(NAME)_$(OS)_$(ARCH)/$(NAME) $(NAME)

dist/$(NAME)_$(OS)_$(ARCH)/$(NAME): $(GOFILES) internal/version.go
	@mkdir -p dist
	goreleaser --snapshot --skip-publish --rm-dist

cmd/$(NAME)/version.go: VERSION
	@test -d cmd/$(NAME) && \
	echo -e '// DO NOT EDIT THIS FILE. Generated from Makefile\npackage main\n\nvar Version = "$(VERSION)"' \
		> cmd/$(NAME)/version.go || exit 0

go.mod:
	@$(DOTENV) make _go.mod

_go.mod:
ifndef GOSERVER
	@make _go.mod_err
else ifndef GOGROUP
	@make _go.mod_err
endif
	go mod init $(GOSERVER)/$(GOGROUP)/$(NAME)
	@make deps

_go.mod_err:
	@echo 'Please run "go mod init server.com/group/project"'
	@echo 'Alternatively set "GOSERVER=$$YOURSERVER" and "GOGROUP=$$YOURGROUP" in ~/.env or $(ROOT)/.env file'
	@exit 1
