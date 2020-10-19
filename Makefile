SHELL = /bin/bash

MAKEFLAGS += --no-print-directory

# Project
NAME := prjstart
GOPATH := $(shell go env GOPATH)
VERSION ?= $(shell cat VERSION)
COMMIT := $(shell test -d .git && git rev-parse --short HEAD)
BUILD_INFO := $(COMMIT)-$(shell date -u +"%Y%m%d-%H%M%SZ")
HASCMD := $(shell test -d cmd && echo "true")
GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?= $(shell uname -m | sed 's/x86_64/amd64/; s/i386/386/')
ARCH = $(shell uname -m)

ISRELEASED := $(shell git show-ref v$$(cat VERSION) 2>&1 > /dev/null && echo "true")

# Utilities
# Default environment variables.
# Any variables already set will override the values in this file(s).
DOTENV := godotenv -f $(HOME)/.env,.env

# Python
PYTHON ?= $(shell command -v python3 python|head -n1)

# Variables
ROOT = $(shell pwd)

# Go
GOMODOPTS = GO111MODULE=on
GOGETOPTS = GO111MODULE=off
GOFILES := $(shell find cmd pkg internal src -name '*.go' 2> /dev/null)
GODIRS = $(shell find . -maxdepth 1 -mindepth 1 -type d | egrep 'cmd|internal|pkg|api')

.PHONY: _build browsereports cattest clean _deps depsdev _go.mod _go.mod_err help \
        _isreleased lint _package _release _release_gitlab test _test _test_setup _test_setup_dirs \
        _test_setup_gitserver _unit _cc _cx _install tag

#
# Help Script
#
define PRINT_HELP_PYSCRIPT
import re, sys

print("Usage: make <target>\n")
cmds = []
for line in sys.stdin:
	match = re.match(r'^_?([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
	  target, help = match.groups()
	  cmds.append([target, help])
for cmd, help in cmds:
		print("  %s%s%s - %s" % ("\x1b[0001m", cmd, "\x1b[0000m", help))
print("")
endef
export PRINT_HELP_PYSCRIPT

#
# End user targets
#
ifneq (, ${PYTHON})
help: ## Print Help
	@$(PYTHON) -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)
else
help:
	$(error python required for 'make help', executable not found)
endif


_build: ## Build binary
	@test -d .cache || go fmt ./...
ifeq ($(XCOMPILE),true)
	GOOS=linux GOARCH=amd64 $(MAKE) dist/$(NAME)_linux_amd64/$(NAME)
	GOOS=darwin GOARCH=amd64 $(MAKE) dist/$(NAME)_darwin_amd64/$(NAME)
	GOOS=windows GOARCH=amd64 $(MAKE) dist/$(NAME)_windows_amd64/$(NAME).exe
endif
ifeq ($(HASCMD),true)
	@$(MAKE) $(NAME)
endif

_install: $(GOPATH)/bin/$(NAME) ## Install to $(GOPATH)/bin

clean: ## Reset project to original state
	-test -f tmp/server.pid && kill -TERM $$(cat tmp/server.pid)
	rm -rf .cache prjstart dist reports tmp vendor

test: ## Test
	$(MAKE) unit
	@exit $$(cat reports/exitcode.txt)

_unit: test_setup ## Unit testing
	@$(MAKE) _test_setup_gitserver
	### Unit Tests
	gotestsum --junitfile reports/junit.xml -- -timeout 5s -covermode atomic -coverprofile=./reports/coverage.out -v ./...; echo $$? > reports/exitcode.txt

_cc: _unit ## Code coverage
	### Code Coverage
	@go tool cover -func=./reports/coverage.out | tee ./reports/coverage.txt
	@go tool cover -html=reports/coverage.out -o reports/html/coverage.html

_cx: test_setup ## Code complexity test
	### Cyclomatix Complexity Report
	@gocyclo -avg $(GODIRS) | grep -v _test.go | tee reports/cyclocomplexity.txt

_package: ## Create an RPM & DEB
	@VERSION=$(VERSION) envsubst < nfpm.yaml.in > nfpm.yaml
	make dist/$(NAME)-$(VERSION).$(ARCH).rpm
	make dist/$(NAME)_$(VERSION)_$(GOARCH).deb

_test_setup: ## Setup test directories
	@mkdir -p tmp
	@mkdir -p reports/html
	@$(MAKE) _test_setup_dirs 2> /dev/null > /dev/null
	@$(MAKE) _test_setup_gitserver 2> /dev/null > /dev/null
	@sync

_test_setup_dirs:
	@cp -r test/fixtures/home tmp/
	@cp -r test/fixtures/checksum tmp/
	@cp -r test/fixtures/compression tmp/
	@cp -r test/fixtures/TestInstall tmp/

_test_setup_gitserver:
	@mkdir -p tmp/gitserveclient
	@-kill -0 $$(cat tmp/server.pid 2>/dev/null) >/dev/null 2>&1 || go run test/fixtures/testserver.go
	@echo "Waiting for git server to launch on 5000..."
	@bash -c 'while ! nc -z localhost 5000; do sleep 0.1; done'
	@echo "git server launched"
	@$(MAKE) _test_setup_gitclient
	@$(MAKE) _test_setup_metadata 2> /dev/null

_test_setup_gitclient:
	@-(find test/fixtures/gitserve -mindepth 1 -maxdepth 1 -type d -exec cp -r {} tmp/gitserveclient \;) 2>&1 > /dev/null
	@-(for i in $$(pwd)/tmp/gitserveclient/*; do cd $$i; git init; git add .; git commit -m "Initial commit"; git tag 7.7.7; git push --set-upstream http://127.0.0.1:5000/$$(basename $$(pwd)).git master; git push --tags; done) 2> /dev/null > /dev/null
	@sync

_test_setup_metadata:
	@-rm -rf tmp/metadata 2> /dev/null > /dev/null
	@mkdir -p tmp/metadata
	@-(find test/fixtures/metadata -mindepth 1 -maxdepth 1 -type d -exec cp -r {} tmp/metadata \;) 2>&1 > /dev/null
	@-(for i in $$(pwd)/tmp/metadata/templates/*; do cd $$i; git init; git add .; git commit -m 'Initial commit';make release; make bumpmajor; make release; git push --set-upstream http://127.0.0.1:5000/$$(basename $$(pwd)).git master; git push --tags; done) 2> /dev/null > /dev/null

_release: ## Release 
	@echo "### Releasing v$(VERSION)"
	@$(MAKE) --no-print-directory _isreleased 2> /dev/null
	git tag v$(VERSION)
	git push --tags

lint: ## Lint tests
	golangci-lint run --enable=gocyclo

tag:
	git fetch --tags
	git tag v$(VERSION)
	git push --tags

deps: go.mod ## Install build dependencies
	$(GOMODOPTS) go mod tidy
	$(GOMODOPTS) go mod vendor

depsdev: ## Install development dependencies
ifeq ($(USEGITLAB),true)
	@mkdir -p $(ROOT)/.cache/{go,gomod}
endif
	@GO111MODULE=on $(MAKE) $(GOGETS)

bumpmajor: ## Version - major bump
	git fetch --tags
	versionbump --checktags major VERSION

bumpminor: ## Version - minor bump
	git fetch --tags
	versionbump --checktags minor VERSION

bumppatch: ## Version - patch bump
	git fetch --tags
	versionbump --checktags patch VERSION

browsereports: _cc ## Open reports in a browser
	@$(MAKE) $(REPORTS)

cattest: ## Print the output of the last set of tests
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
$(GOPATH)/bin/$(NAME): $(NAME)
	install -m 755 $(NAME) $(GOPATH)/bin/$(NAME)

GOGETS := github.com/jstemmer/go-junit-report gotest.tools/gotestsum github.com/golangci/golangci-lint/cmd/golangci-lint@v1.25.0 \
	      github.com/goreleaser/goreleaser github.com/fzipp/gocyclo github.com/joho/godotenv/cmd/godotenv \
		  github.com/crosseyed/versionbump/cmd/versionbump github.com/sosedoff/gitkit
.PHONY: $(GOGETS)
$(GOGETS):
	cd /tmp; go get $@

REPORTS = reports/html/coverage.html
.PHONY: $(REPORTS)
$(REPORTS):
ifeq ($(GOOS),darwin)
	@test -f $@ && open $@
else ifeq ($(GOOS),linux)
	@test -f $@ && xdg-open $@
endif

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
$(NAME): dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	install -m 755 $< $@

dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME) dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME).exe: $(GOFILES) internal/version.go
	@mkdir -p $$(dirname $@)
	go build -o $@ ./cmd/prjstart

dist/$(NAME)-$(VERSION).$(ARCH).rpm: dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	@mkdir -p $$(dirname $@)
	@make nfpm.yaml
	nfpm pkg --packager rpm --target dist/

dist/$(NAME)_$(VERSION)_$(GOARCH).deb: dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	@mkdir -p $$(dirname $@)
	@make nfpm.yaml
	nfpm pkg --packager deb --target dist/

internal/version.go: VERSION
	@test -d $$(dirname $@) && \
	echo -e '// DO NOT EDIT THIS FILE. Makefile generated.\n//\npackage internal\n\nvar Version = "$(VERSION)"' \
		> $@ || exit 0

nfpm.yaml: nfpm.yaml.in VERSION
	@VERSION=$(VERSION) $(DOTENV) envsubst < $< > $@

go.mod:
	@$(DOTENV) $(MAKE) _go.mod

_go.mod:
ifndef GOSERVER
	@$(MAKE) _go.mod_err
else ifndef GOGROUP
	@$(MAKE) _go.mod_err
endif
	go mod init $(GOSERVER)/$(GOGROUP)/$(NAME)
	@$(MAKE) deps

_go.mod_err:
	@echo 'Please run "go mod init server.com/group/project"'
	@echo 'Alternatively set "GOSERVER=$$YOURSERVER" and "GOGROUP=$$YOURGROUP" in ~/.env or $(ROOT)/.env file'
	@exit 1

#
# make wrapper - Execute any target target prefixed with a underscore.
# EG 'make vmcreate' will result in the execution of 'make _vmcreate' 
#
%:
	@egrep -q '^_$@:' Makefile && $(DOTENV) $(MAKE) _$@
