VERSION := $(shell git describe --tags 2>/dev/null)
DIST_DIRS := find * -type d -exec
export GO15VENDOREXPERIMENT=1

build:
	go build -o helm.bin -ldflags "-X main.version=${VERSION}" helm/helm.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./helm.bin ${DESTDIR}/usr/local/bin/helm

test:
	go test ./helm/. ./helm/manifest ./helm/action ./helm/log ./helm/model

clean:
	rm -f ./helm/helm.test
	rm -f ./helm

bootstrap:
	glide up

bootstrap-dist:
	go get -u github.com/mitchellh/gox

build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows " \
	-arch="amd64 386" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE.txt {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) zip -r helm-{}.zip {} \; && \
	cd ..

test-charts:
	@./test/test-charts $(TEST_CHARTS)

.PHONY: build test install clean bootstrap bootstrap-dist build-all dist test-charts
