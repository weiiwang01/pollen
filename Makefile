GO_BUILD=go build
GO_TEST=go test
GO_CLEAN=go clean

all: pollen

pollen: pollen.go metrics.go
	$(GO_BUILD) -o $@ $^

test: pollen.go pollen_test.go metrics.go metrics_test.go
	$(GO_TEST)

clean:
	$(RM) pollen

.PHONY: all clean test
