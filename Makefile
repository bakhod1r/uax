.PHONY: all test race bench compare fuzz lint cover clean

all: lint test

test:
	go test ./...
	@for m in middleware/gin middleware/echo middleware/fiber; do \
		echo "--- $$m"; (cd $$m && go build ./... && go test ./...) || exit 1; \
	done

race:
	go test -race ./...

# The allocation budget is an assertion, not a report.
bench:
	go test -run=XXX -bench=. -benchmem .

compare:
	@cd benchmark && go test -run=XXX -bench=. -benchmem .

fuzz:
	go test -run=XXX -fuzz=FuzzParse -fuzztime=60s .
	go test -run=XXX -fuzz=FuzzScan -fuzztime=60s ./internal/scan

lint:
	gofmt -l .
	go vet ./...

cover:
	go test -coverprofile=/tmp/uax.cover ./... && go tool cover -func=/tmp/uax.cover | tail -1

clean:
	go clean -testcache
