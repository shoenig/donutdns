NAME = donutdns

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o output/$(NAME)

.PHONY: clean
clean:
	rm -rf dist output/$(NAME)

.PHONY: test
test:
	go test -race ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: release
release:
	envy exec gh-release goreleaser release --clean
	$(MAKE) clean

default: build

