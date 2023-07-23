NAME = donutdns

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -o $(NAME)

.PHONY: clean
clean:
	rm -rf dist $(NAME)

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

