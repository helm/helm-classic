ifndef GOPATH
$(error No GOPATH set)
endif

BIN_DIR := bin
DIST_DIR := _dist
GO_PACKAGES := action chart config dependency log manifest release plugins/sec plugins/example codec
MAIN_GO := helm.go
HELM_BIN := $(BIN_DIR)/helm
PATH_WITH_HELM = PATH="$(shell pwd)/$(BIN_DIR):$(PATH)"

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)+$(shell git rev-parse --short HEAD)

export GO15VENDOREXPERIMENT=1

ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

# dockerized development environment variables
REPO_PATH := github.com/helm/helm
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}


build: $(MAIN_GO)
	${DEV_ENV_CMD} go build -o $(HELM_BIN) -ldflags "-X github.com/helm/helm/cli.version=${VERSION}" $<

bootstrap:
	${DEV_ENV_CMD} glide install

build-all:
	${DEV_ENV_CMD} gox -verbose \
		-ldflags "-X github.com/helm/helm/cli.version=${VERSION}" \
		-os="linux darwin " \
		-arch="amd64 386" \
		-output="$(DIST_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

clean:
	rm -rf $(DIST_DIR) $(BIN_DIR)

dist: build-all
	@cd $(DIST_DIR) && \
	find * -type d -exec zip -jr helm-$(VERSION)-{}.zip {} \; && \
	cd -

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 $(HELM_BIN) ${DESTDIR}/usr/local/bin/helm

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
	$(PATH_WITH_HELM) go test -short ./ $(addprefix ./,$(GO_PACKAGES))

test: test-style
	$(PATH_WITH_HELM) go test -v ./ $(addprefix ./,$(GO_PACKAGES))

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
