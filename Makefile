GO_BUILD=go build
GO_TEST=go test
GO_CLEAN=go clean
GIT_ARCHIVE=git archive

VERSION=4.23
TAG=v4.23

all: pollen

pollen: pollen.go metrics.go
	$(GO_BUILD) -o $@ $^

test: pollen.go pollen_test.go metrics.go metrics_test.go
	$(GO_TEST)

dist: pollen
	git tag $(TAG)
	$(GIT_ARCHIVE) --format=tar --prefix=pollen-$(VERSION)/ $(TAG) | bzip2 > ../pollen-$(VERSION).tar.bz2

clean:
	$(RM) pollen

snap: clean
	export SNAPCRAFT_BUILD_ENVIRONMENT_MEMORY=4G
	export SNAPCRAFT_BUILD_ENVIRONMENT_CPU=4
	export SNAPCRAFT_MAX_PARALLEL_BUILD_COUNT=4
	snapcraft

.PHONY: all clean test
