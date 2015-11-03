VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)+$(shell git rev-parse --short HEAD)
DIST_DIRS := find * -type d -exec
export GO15VENDOREXPERIMENT=1

ifndef VERSION
  VERSION := git-$(shell git rev-parse --short HEAD)
endif

bootstrap:
	glide -y glide-full.yaml up

bootstrap-dist:
	go get -u github.com/mitchellh/gox

build:
	go build -o bin/helm -ldflags "-X main.version=${VERSION}" helm.go

build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin " \
	-arch="amd64 386" \
	-output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" . && \
	cd ..

clean:
	rm -f ./bin/helm

dist: build-all
	@mkdir -p _dist
	@cd _dist && \
	$(DIST_DIRS) zip -jr helm-$(VERSION)-{}.zip {} \; && \
	cd ..

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 bin/helm ${DESTDIR}/usr/local/bin/helm

prep-bintray-json:
# TRAVIS_TAG is set to the tag name if the build is a tag
ifdef TRAVIS_TAG
	@jq '.version.name |= "$(VERSION)"' ci/bintray-template.json | \
		jq '.package.repo |= "helm"' > ci/bintray-ci.json
else
	@jq '.version.name |= "$(VERSION)"' ci/bintray-template.json \
		> ci/bintray-ci.json
endif

quicktest:
	go test ./. ./manifest ./action ./log ./chart ./dependency ./repo ./release

test:
	go test -v ./. ./manifest ./action ./log ./chart ./dependency ./repo ./release

test-charts:
	@./_test/test-charts $(TEST_CHARTS)

.PHONY: bootstrap \
				bootstrap-dist \
				build \
				build-all \
				clean \
				dist \
				install \
				prep-bintray-json \
				test \
				test-charts
