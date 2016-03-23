SHORT_NAME := helm
DEIS_REGISTRY ?= ${DEV_REGISTRY}
IMAGE_PREFIX ?= helm

REPO_PATH := github.com/helm/${SHORT_NAME}

# The following variables describe the containerized development environment
# and other build options
BIN_DIR := bin
DIST_DIR := _dist
GO_PACKAGES := action chart config dependency log manifest release plugins/sec plugins/example codec
MAIN_GO := helm.go
HELM_BIN := $(BIN_DIR)/helm
PATH_WITH_HELM = PATH="$(shell pwd)/$(BIN_DIR):$(PATH)"
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)+$(shell git rev-parse --short HEAD)
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
DEV_ENV_CMD_INT := docker run -it --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}
LDFLAGS := "-s -X main.version=${VERSION}"



# Allow developers to step into the containerized development environment

# Containerized dependency resolution

bootstrap: check-docker 
	${DEV_ENV_CMD} glide install

check-docker:
	@if [ -z $$(which docker) ]; then \
		echo "Missing \`docker\` client which is required for development"; \
		exit 2; \
	fi
# Containerized build of the binary
build: check-docker
	mkdir -p ${BIN_DIR}
	make binary-build

docker-build: build check-docker
	docker build --rm -t ${IMAGE} rootfs
	docker tag -f ${IMAGE} ${MUTABLE_IMAGE}


# Builds the binary-- this should only be executed within the
# containerized development environment.

binary-build:
	make build-all
	@cd $(DIST_DIR) && \
	find * -type d -exec zip -jr helm-$(VERSION)-{}.zip {} \; && \
	cd -

build-all:
	${DEV_ENV_CMD} gox -verbose \
	-ldflags "-X github.com/helm/helm/cli.version=${VERSION}" \
	-os="linux darwin " \
	-arch="amd64 386" \
	-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

prep-bintray-json:
# TRAVIS_TAG is set to the tag name if the build is a tag
ifdef TRAVIS_TAG
	@jq '.version.name |= "$(VERSION)"' _scripts/ci/bintray-template.json | \
		jq '.package.repo |= "helm"' > _scripts/ci/bintray-ci.json
else
	@jq '.version.name |= "$(VERSION)"' _scripts/ci/bintray-template.json \
		> _scripts/ci/bintray-ci.json
endif

quicktest:
	${DEV_ENV_CMD} go test -short ./ $(addprefix ./,$(GO_PACKAGES))

test:
	${DEV_ENV_CMD} go test -v ./ $(addprefix ./,$(GO_PACKAGES))

test-style:
	@if [ $(shell gofmt -e -l -s *.go $(GO_PACKAGES)) ]; then \
		echo "gofmt check failed:"; gofmt -e -l -s *.go $(GO_PACKAGES); exit 1; \
	fi
	@for i in . $(GO_PACKAGES); do \
		golint $$i; \
	done
	@for i in . $(GO_PACKAGES); do \
		go vet github.com/helm/helm/$$i; \
	done

clean: check-docker
	docker rmi ${IMAGE}


.PHONY: bootstrap \
				build \
				build-all \
				clean \
				dist \
				install \
				prep-bintray-json \
				quicktest \
				test \
				test-charts \
				test-style
