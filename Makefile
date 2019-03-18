
.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	# Ignore grep's exit code since no match returns 1.
	echo 'linting...' ; golint ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: all
all: vet lint test

.PHONY: example
