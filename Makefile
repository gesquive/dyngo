#
#  Makefile
#
#  The kickoff point for all project management commands.
#

# TODO: Add cronjob install
# TODO: Add service install
GOCC := go

# Program version
VERSION := $(shell git describe --always --tags)

# Binary name for bintray
BIN_NAME=digitalocean-ddns

# Project owner for bintray
OWNER=gesquive

# Project name for bintray
PROJECT_NAME=digitalocean-ddns

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com

# Grab the current commit
GIT_COMMIT=$(shell git rev-parse HEAD)

# Check if there are uncommited changes
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Use a local vendor directory for any dependencies; comment this out to
# use the global GOPATH instead
# GOPATH=$(PWD)

INSTALL_PATH=$(GOPATH)/${REPO_HOST_URL}/${OWNER}/${PROJECT_NAME}

FIND_DIST:=find * -type d -exec

default: build

help:
	@echo 'Management commands for $(BIN_NAME):'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Compile the project
	@echo "building ${OWNER} ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	${GOCC} build -ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" -o ${BIN_NAME}

install: build ## Install binary
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./${BIN_NAME} ${DESTDIR}/usr/local/bin/${BIN_NAME}

depends: ## Download project dependencies
	${GOCC} get -u github.com/Masterminds/glide
	glide install

test: ## Run golang tests
	${GOCC} test ./...

clean: ## Clean the directory tree
	${GOCC} clean
	rm -f ./${BIN_NAME}.test
	rm -f ./${BIN_NAME}
	rm -rf ./dist

bootstrap-dist:
	${GOCC} get -u github.com/mitchellh/gox

build-all: bootstrap-dist
	gox -verbose \
	-ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" \
	-os="linux darwin windows " \
	-arch="amd64 386" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all
	cd dist && \
	$(FIND_DIST) cp ../LICENSE {} \; && \
	$(FIND_DIST) cp ../README.md {} \; && \
	$(FIND_DIST) tar -zcf ${PROJECT_NAME}-${VERSION}-{}.tar.gz {} \; && \
	$(FIND_DIST) zip -r ${PROJECT_NAME}-${VERSION}-{}.zip {} \; && \
	cd ..

fmt: ## Reformat the source tree with gofmt
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

link: ## Symlink this project into the GOPATH
	# relink into the go path
	if [ ! $(INSTALL_PATH) -ef . ]; then \
		mkdir -p `dirname $(INSTALL_PATH)`; \
		ln -s $(PWD) $(INSTALL_PATH); \
	fi


.PHONY: build help test install depends clean bootstrap-dist build-all dist fmt link
