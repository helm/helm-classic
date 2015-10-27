VERSION := $(shell git describe --tags 2>/dev/null)
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
	go build -o bin/helm -ldflags "-X main.version=${VERSION}" helm/helm.go

build-all:
	@cd helm && \
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin " \
	-arch="amd64 386" \
	-output="../dist/{{.OS}}-{{.Arch}}/{{.Dir}}" . && \
	cd ..

clean:
	rm -f ./helm/helm.test
	rm -f ./helm.bin

dist: build-all
	@cd dist && \
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
	go test ./helm/. ./helm/manifest ./helm/action ./helm/log ./helm/model ./helm/dependency

test:
	go test -v ./helm/. ./helm/manifest ./helm/action ./helm/log ./helm/model ./helm/dependency

test-charts:
	@./test/test-charts $(TEST_CHARTS)

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
