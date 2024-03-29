SHELL = /bin/bash

MAKEFLAGS += --no-print-directory

# Project
NAME := kick
GOPATH := $(shell go env GOPATH)
VERSION ?= $(shell cat VERSION)
HASCMD := $(shell test -d cmd && echo "true")
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
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

#
# End user targets
#

### HELP

.PHONY: help
help: ## Print Help
	@maker --menu=Makefile

### DEVELOPMENT

.PHONY: _build
_build: ## Build binary
	@test -d .cache || go fmt ./...
ifeq ($(HASCMD),true)
ifeq ($(XCOMPILE),true)
	GOOS=linux GOARCH=amd64 $(MAKE) dist/$(NAME)_linux_amd64/$(NAME)
	GOOS=darwin GOARCH=amd64 $(MAKE) dist/$(NAME)_darwin_amd64/$(NAME)
	GOOS=windows GOARCH=amd64 $(MAKE) dist/$(NAME)_windows_amd64/$(NAME).exe
endif
	@$(MAKE) $(NAME)
endif

.PHONY: _install
_install: $(GOPATH)/bin/$(NAME) ## Install to $(GOPATH)/bin

.PHONY: clean
clean: ## Reset project to original state
	-test -f tmp/server-$(USER).pid && kill -TERM $$(cat tmp/server-$(USER).pid)
	rm -rf .cache kick dist reports tmp vendor nfpm.yaml www/site

.PHONY: test
test: ## Test
	$(MAKE) test_setup
	$(MAKE) lint
	@#echo FEATURE FLAGS ENABLED
	@#FF_ENABLED=true $(MAKE) _test_fflags
	@echo FEATURE FLAGS DISABLED
	FF_ENABLED=false $(MAKE) _test_std
	@# Combined the return codes of all the tests
	@echo "Exit codes, unit tests: $$(cat reports/exitcode-unit.txt), golangci-lint: $$(cat reports/exitcode-golangci-lint.txt), golint: $$(cat reports/exitcode-golint.txt)"
	@exit $$(( $$(cat reports/exitcode-unit.txt) + $$(cat reports/exitcode-golangci-lint.txt) + $$(cat reports/exitcode-golint.txt) ))

# Test without feature flags
.PHONY: _test_std
_test_std:
	$(MAKE) unit
	$(MAKE) cx
	$(MAKE) cc

# Test with feature flags
.PHONY: _test_fflags
_test_fflags:
	$(MAKE) unit

.PHONY: _unit
_unit:
	### Unit Tests
	gotestsum --jsonfile reports/unit.json --junitfile reports/junit.xml -- -timeout 60s -covermode atomic -coverprofile=./reports/coverage.out -v ./...; echo $$? > reports/exitcode-unit.txt
	@go-test-report -t "kick unit tests" -o reports/html/unit.html < reports/unit.json > /dev/null

.PHONY: _cc
_cc:
	### Code Coverage
	@go-acc -o ./reports/coverage.out ./... > /dev/null
	@go tool cover -func=./reports/coverage.out | tee reports/coverage.txt
	@go tool cover -html=reports/coverage.out -o reports/html/coverage.html

.PHONY: _cx
_cx:
	### Cyclomatix Complexity Report
	@gocyclo -avg $(GODIRS) | grep -v _test.go | tee reports/cyclomaticcomplexity.txt
	@contents=$$(cat reports/cyclomaticcomplexity.txt); echo "<html><title>cyclomatic complexity</title><body><pre>$${contents}</pre></body><html>" > reports/html/cyclomaticcomplexity.html

.PHONY: _package
_package: ## Create an RPM, Deb, Homebrew package
	@XCOMPILE=true make build
	@VERSION=$(VERSION) envsubst < nfpm.yaml.in > nfpm.yaml
	$(MAKE) dist/kick-$(VERSION).rb
	$(MAKE) tmp/kick.rb
	$(MAKE) dist/$(NAME)-$(VERSION).$(ARCH).rpm
	$(MAKE) dist/$(NAME)_$(VERSION)_$(GOARCH).deb

.PHONY: _test_setup
_test_setup:
	@mkdir -p tmp
	@mkdir -p reports/html
	@$(MAKE) _test_setup_dirs 2> /dev/null > /dev/null
	@$(MAKE) _test_setup_gitserver 2> /dev/null > /dev/null
	@sync

.PHONY: _test_setup_dirs
_test_setup_dirs:
	@find test/fixtures -mindepth 1 -maxdepth 1 -type d | grep -v gitserve | xargs -I {} cp -r {} tmp/
	@rm -rf tmp/TestInfo
	@mkdir -p tmp/TestInfo/testrepo
	@touch tmp/TestInfo/testrepo/empty
	@cd tmp/TestInfo/testrepo; (git init -q; git add .; git commit -m 'Initial commit'; git tag 1.0.0; git tag 1.1.0; git tag 2.0.0; git tag 2.1.0; git tag 2.1.1) > /dev/null

.PHONY: _test_setup_gitserver
_test_setup_gitserver:
	-kill -TERM $$(cat tmp/server-$(USER).pid 2>/dev/null) >/dev/null 2>&1
	rm -rf tmp/gitserve 2> /dev/null > /dev/null
	set -e; find test/fixtures/gitserve -mindepth 1 -maxdepth 1 -type d | xargs -I {} basename {} | xargs -I {} bash -c "set -e; mkdir -p tmp/gitserve/{}.git; cd tmp/gitserve/{}.git; git init --bare"
	go run test/fixtures/testserver.go -p tmp/server-$(USER).pid
	mkdir -p tmp/gitserveclient
	echo "Waiting for git server to launch on 8080..."
	bash -c 'while ! nc -z localhost 8080; do sleep 0.1; done'
	echo "git server launched"
	$(MAKE) _test_setup_gitclient

.PHONY: _test_setup_gitclient
_test_setup_gitclient:
	rm -rf tmp/gitserverclient 2> /dev/null > /dev/null
	(set -e; find test/fixtures/gitserve -mindepth 1 -maxdepth 1 -type d -exec cp -r {} tmp/gitserveclient \;) 2>&1 > /dev/null
	find tmp/gitserveclient -type d -mindepth 1 -maxdepth 1 | xargs -I {} -n 1 bash -c "cd {}; git init; git add .; git commit -m 'Initial commit'; git tag 7.7.7; git remote add origin http://127.0.0.1:8080/\$$(basename {}).git; git push -u origin master; git push --tags"
	sync

.PHONY: _release
_release: ## Trigger a release by creating a tag and pushing to the upstream repository
	@echo "### Releasing v$(VERSION)"
	@$(MAKE) _isreleased 2> /dev/null
	git tag v$(VERSION)
	git push --tags
	$(MAKE) deploydocs

# To be run inside a github workflow
.PHONY: _release_github
_release_github: _package
	gh release create v$(VERSION)
	gh release upload v$(VERSION) dist/kick-$(VERSION).tar.gz
	gh release upload v$(VERSION) dist/kick.rb
	gh release upload v$(VERSION) dist/kick-$(VERSION).rb
	gh release upload v$(VERSION) dist/kick-$(VERSION).x86_64.rpm
	gh release upload v$(VERSION) dist/kick_$(VERSION)_amd64.deb

.PHONY: lint
lint: internal/version.go
	golangci-lint run --enable=gocyclo --timeout 120s; echo $$? > reports/exitcode-golangci-lint.txt
	golint -set_exit_status ./..; echo $$? > reports/exitcode-golint.txt

.PHONY: tag
tag:
	git fetch --tags
	git tag v$(VERSION)
	git push --tags

.PHONY: deps
deps: go.mod ## Install build dependencies
	$(GOMODOPTS) go mod tidy
	$(GOMODOPTS) go mod download

.PHONY: depsdev
depsdev: ## Install development dependencies
ifeq ($(USEGITLAB),true)
	@mkdir -p $(ROOT)/.cache/{go,gomod}
endif
	cat goinstalls.txt | egrep -v '^#' | xargs -t -n1 go install

.PHONY: report
report: ## Open reports generated by "make test" in a browser
	@$(MAKE) $(REPORTS)

### VERSION INCREMENT

.PHONY: bumpmajor
bumpmajor: ## Increment VERSION file ${major}.0.0 - major bump
	git fetch --tags
	versionbump --checktags major VERSION

.PHONY: bumpminor
bumpminor: ## Increment VERSION file 0.${minor}.0 - minor bump
	git fetch --tags
	versionbump --checktags minor VERSION

.PHONY: bumppatch
bumppatch: ## Increment VERSION file 0.0.${patch} - patch bump
	git fetch --tags
	versionbump --checktags patch VERSION

.PHONY: getversion
getversion:
	VERSION=$(VERSION) bash -c 'echo $$VERSION'

### DOCUMENTATION

.PHONY: builddocs
builddocs: www/docs/cli.md www/docs/install.md ## Build documents

deploydocs: builddocs ## Deploy documentaiton to github
	cd www; mkdocs gh-deploy

.PHONY: _docserver
_docserver: builddocs ## Start document server
	cd www; mkdocs serve --dev-addr localhost:$(KICK_MKDOCS_PORT)

#
# Helper targets
#
$(GOPATH)/bin/$(NAME): $(NAME)
	install -m 755 $(NAME) $(GOPATH)/bin/$(NAME)

# Open html reports
REPORTS = reports/html/unit.html reports/html/coverage.html reports/html/cyclomaticcomplexity.html
.PHONY: $(REPORTS)
$(REPORTS):
ifeq ($(GOOS),darwin)
	@test -f $@ && open $@
else ifeq ($(GOOS),linux)
	@test -f $@ && xdg-open $@
endif

# Check versionbump
.PHONY: _isreleased
_isreleased:
ifeq ($(ISRELEASED),true)
	@echo "Version $(VERSION) has been released."
	@echo "Please bump with 'make bump(minor|patch|major)' depending on breaking changes."
	@exit 1
endif

### VAGRANT

.PHONY: up
up: ## Start vagrant
	@(vagrant status --no-color | grep running 2>&1) > /dev/null || VAGRANT_EXPERIMENTAL=disks vagrant up

.PHONY: ssh
ssh: up ## Run vagrant ssh and cd to shared directory
	vagrant ssh -c 'cd /vagrant; exec $$SHELL -l'

.PHONY: halt
halt: ## Run vagrant halt
	vagrant halt

.PHONY: destroy
destroy: ## Run vagrant destroy
	vagrant destroy

#
# File targets
#
$(NAME): dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	install -m 755 $< $@

dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME) dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME).exe: $(GOFILES) internal/version.go
	@mkdir -p $$(dirname $@)
	go generate ./...
	go build -tags sqlite_foreign_keys -o $@ ./cmd/kick

dist/$(NAME)-$(VERSION).$(ARCH).rpm: dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	@mkdir -p $$(dirname $@)
	@$(MAKE) nfpm.yaml
	nfpm pkg --packager rpm --target dist/

dist/$(NAME)_$(VERSION)_$(GOARCH).deb: dist/$(NAME)_$(GOOS)_$(GOARCH)/$(NAME)
	@mkdir -p $$(dirname $@)
	@$(MAKE) nfpm.yaml
	nfpm pkg --packager deb --target dist/

internal/version.go: internal/version.go.in VERSION
	@VERSION=$(VERSION) $(DOTENV) envsubst < $< > $@

dist/kick-$(VERSION).rb: dist/kick.rb
	@cp dist/kick{,-$(VERSION)}.rb

dist/kick.rb: kick.rb.in dist/kick-$(VERSION).tar.gz
	@BASEURL="https://github.com/kick-project/kick/releases/download/v$(VERSION)" VERSION=$(VERSION) SHA256=$$(sha256sum dist/kick-$(VERSION).tar.gz | awk '{print $$1}') $(DOTENV) envsubst < $< > $@

tmp/kick.rb: kick.rb.in dist/kick-$(VERSION).tar.gz
	@mkdir -p tmp
	@BASEURL="file://$(PWD)/dist" VERSION=$(VERSION) SHA256=$$(sha256sum dist/kick-$(VERSION).tar.gz | awk '{print $$1}') $(DOTENV) envsubst < $< > $@

nfpm.yaml: nfpm.yaml.in VERSION
	@VERSION=$(VERSION) $(DOTENV) envsubst < $< > $@

dist/kick-$(VERSION).tar.gz: $(GOFILES)
	tar -zcf dist/kick-$(VERSION).tar.gz $$(find . \( -path ./test -prune -o -path ./tmp \) -prune -false -o \( -name go.mod -o -name go.sum -o -name \*.go \))

www/docs/cli.md: www/docs/cli.md.sh $(NAME)
	bash $< > $@

www/docs/install.md: www/docs/install.md.in VERSION
	VERSION=$(VERSION) $(DOTENV) envsubst < $< > $@

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
	@maker $@

.PHONY: _env
_env:
	@echo "DEVELOPMENT:"
	@echo "    XCOMPILE=$(XCOMPILE)"
	@echo "    KICK_DEBUG=$(XDEBUG)"
	@echo "TESTING:"
	@echo "    KICK_LISTEN=$(KICK_LISTEN)"
	@echo "    KICK_TEST_SSH=$(KICK_TEST_SSH)"
	@echo "FEATURE FLAGS"
	@echo "    FF_ENABLED=$(FF_ENABLED)"
