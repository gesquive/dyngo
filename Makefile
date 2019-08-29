#
#  Makefile
#
#  A kickass golang v1.12.x makefile
#  v1.0.1

GOCC := go

# Program version
MK_VERSION := $(shell git describe --always --tags)

# Check if there are uncommited changes
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

PKG_NAME := ${REPO_HOST_URL}/${OWNER}/${PROJECT_NAME}
INSTALL_PATH := ${GOPATH}/src/${PKG_NAME}

DIST_OS ?= "linux darwin windows"
DIST_ARCH ?= "amd64 386"
DIST_ARCHIVE ?= "tar.gz"
DIST_FILES ?= "LICENSE README.md"

COVER_PATH := coverage
DIST_PATH ?= dist
INSTALL_PATH ?= "/usr/local/bin"
PKG_LIST := ./...

BIN ?= ${GOPATH}/bin
GOLINT ?= ${BIN}/golint
GOX ?= ${BIN}/gox
GOP ?= ${BIN}/gop

export SHELL ?= /bin/bash

include make.cfg
default: test build

.PHONY: help
help:
	@echo 'Management commands for $(PROJECT_NAME):'
	@grep -Eh '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	 awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compile the project
	@echo "building ${OWNER} ${BIN_NAME} ${MK_VERSION}"
	@echo "GOPATH=${GOPATH}"
	${GOCC} build -ldflags "-X main.version=${MK_VERSION} -X main.dirty=${GIT_DIRTY}" -o ${BIN_NAME}

.PHONY: install
install: build ## Install the binary
	install -d ${DESTDIR}
	install -m 755 ./${BIN_NAME} ${DESTDIR}/${BIN_NAME}

.PHONY: link
link: $(INSTALL_PATH) ## Symlink this project into the GOPATH
$(INSTALL_PATH):
	@mkdir -p `dirname $(INSTALL_PATH)`
	@ln -s $(PWD) $(INSTALL_PATH) >/dev/null 2>&1

.PHONY: path # Returns the project path
path:
	@echo $(INSTALL_PATH)

.PHONY: deps
deps: ## Download project dependencies
	${GOCC} mod download

.PHONY: lint
lint: ${GOLINT} ## Lint the source code
	${GOLINT} -set_exit_status ${PKG_LIST}

.PHONY: test
test: ## Run golang tests
	${GOCC} test ${PKG_LIST}

.PHONY: bench
bench: ## Run golang benchmarks
	${GOCC} test -benchmem -bench=. ${PKG_LIST}

.PHONY: coverage
coverage: ## Run coverage report
	${GOCC} test -v -cover ${PKG_LIST}

.PHONY: coverage-report
coverage-report: ## Generate global code coverage report
	mkdir -p "${COVER_PATH}"
	${GOCC} test -v -coverprofile "${COVER_PATH}/coverage.dat" ${PKG_LIST}
	${GOCC} tool cover -html="${COVER_PATH}/coverage.dat" -o "${COVER_PATH}/coverage.html"

.PHONY: race
race: ## Run data race detector
	${GOCC} test -race ${PKG_LIST}

.PHONY: clean
clean: ## Clean the directory tree
	${GOCC} clean
	rm -f ./${BIN_NAME}.test
	rm -f ./${BIN_NAME}
	rm -rf "${DIST_PATH}"
	rm -f "${COVER_PATH}"

.PHONY: build-dist
build-dist: ${GOX}
	${GOX} -verbose \
	-ldflags "-X main.version=${MK_VERSION} -X main.dirty=${GIT_DIRTY}" \
	-os=${DIST_OS} \
	-arch=${DIST_ARCH} \
	-output="${DIST_PATH}/{{.OS}}-{{.Arch}}/{{.Dir}}" .

.PHONY: package-dist
package-dist: ${GOP}
	${GOP} --delete \
	--os=${DIST_OS} \
	--arch=${DIST_ARCH} \
	--archive=${DIST_ARCHIVE} \
	--files=${DIST_FILES} \
	--input="${DIST_PATH}/{{.OS}}-{{.Arch}}/{{.Dir}}" \
	--output="${DIST_PATH}/{{.Dir}}-${MK_VERSION}-{{.OS}}-{{.Arch}}.{{.Archive}}" .

.PHONY: dist
dist: build-dist package-dist ## Cross compile and package the full distribution

.PHONY: fmt
fmt: ## Reformat the source tree with gofmt
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

# Install golang dependencies here
${BIN}:
	@echo "Making bin"
	@mkdir -p $@
${BIN}/%: ${BIN}
	@echo "Installing ${PACKAGE} to ${BIN}"
	@tmp=$$(mktemp -d); \
       env GO111MODULE=off GOPATH=$$tmp GOBIN=${BIN} ${GOCC} get ${PACKAGE} \
        || ret=$$?; \
       rm -rf $$tmp ; exit $$ret

${BIN}/golint: PACKAGE=golang.org/x/lint/golint
${BIN}/gox:    PACKAGE=github.com/mitchellh/gox
${BIN}/gop:    PACKAGE=github.com/gesquive/gop
