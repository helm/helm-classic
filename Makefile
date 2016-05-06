export GO15VENDOREXPERIMENT=1

REPO_PATH := github.com/helm/helm-classic

# The following variables describe the containerized development environment
# and other build options
BIN_DIR := bin
DIST_DIR := _dist
GO_PACKAGES := action chart config dependency log manifest release plugins/sec plugins/example codec
MAIN_GO := helmc.go
HELM_BIN := ${BIN_DIR}/helmc

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)+$(shell git rev-parse --short HEAD)

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.12.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
DEV_ENV_CMD_INT := docker run -it --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
LDFLAGS := "-X ${REPO_PATH}/cli.version=${VERSION}"

PATH_WITH_HELM = PATH=${DEV_ENV_WORK_DIR}/${BIN_DIR}:$$PATH

check-docker:
	@if [ -z $$(which docker) ]; then \
		echo "Missing \`docker\` client which is required for development and testing"; \
		exit 2; \
	fi

# Allow developers to step into the containerized development environment
dev: check-docker
	${DEV_ENV_CMD_INT} bash

# Containerized dependency resolution
bootstrap: check-docker
	${DEV_ENV_CMD} glide install

# Containerized build of the binary
build: check-docker
	${DEV_ENV_CMD} make native-build

# Builds the binary for the native OS and architecture.
# This can be run directly to compile for one's own system if desired.
# It can also be run within a container (which will compile for Linux/64) using `make build`
native-build:
	go build -o ${HELM_BIN} -ldflags ${LDFLAGS} ${MAIN_GO}

# Containerized build of binaries for all supported OS and architectures
build-all: check-docker
	${DEV_ENV_CMD} gox -verbose \
	-ldflags ${LDFLAGS} \
	-os="linux darwin " \
	-arch="amd64 386" \
	-output="${DIST_DIR}/{{.OS}}-{{.Arch}}/helmc" .

clean:
	rm -rf ${DIST_DIR} ${BIN_DIR}

dist: build-all
	${DEV_ENV_CMD} bash -c 'cd ${DIST_DIR} && find * -type d -exec zip -jr helmc-${VERSION}-{}.zip {} \;'

install:
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ${HELM_BIN} ${DESTDIR}/usr/local/bin/helmc

prep-bintray-json:
# TRAVIS_TAG is set to the tag name if the build is a tag
ifdef TRAVIS_TAG
	${DEV_ENV_CMD} jq '.version.name |= "${VERSION}"' _scripts/ci/bintray-template.json | \
		jq '.package.repo |= "helm"' > _scripts/ci/bintray-ci.json
else
	${DEV_ENV_CMD} jq '.version.name |= "${VERSION}"' _scripts/ci/bintray-template.json \
		> _scripts/ci/bintray-ci.json
endif

quicktest:
	${DEV_ENV_CMD} bash -c '${PATH_WITH_HELM} go test -short ./ $(addprefix ./,${GO_PACKAGES})'

test: test-style
	${DEV_ENV_CMD} bash -c '${PATH_WITH_HELM} go test -v ./ $(addprefix ./,${GO_PACKAGES})'

test-style:
	${DEV_ENV_CMD} gofmt -e -l -s *.go ${GO_PACKAGES}
	@${DEV_ENV_CMD} bash -c 'gofmt -e -l -s *.go ${GO_PACKAGES} | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi'
	@${DEV_ENV_CMD} bash -c 'for i in . ${GO_PACKAGES}; do golint $$i; done'
	@${DEV_ENV_CMD} bash -c 'for i in . ${GO_PACKAGES}; do go vet ${REPO_PATH}/$$i; done'

.PHONY: check-docker \
				dev \
				bootstrap \
				build \
				native-build \
				build-all \
				clean \
				dist \
				install \
				prep-bintray-json \
				quicktest \
				test \
				test-style
