GO_BUILD=go build
GO_TEST=go test
GO_CLEAN=go clean
GIT_ARCHIVE=git archive

VERSION=4.22
TAG=v4.22

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

.PHONY: all clean test
