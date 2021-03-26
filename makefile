all: build

ENVVAR = GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
TAG = v0.0.1

build: clean
	$(ENVVAR) go build -o cloud-platform-applier

container:
	docker build -t cloud-platform-applier:$(TAG) .

clean:
	rm -f cloud-platform-applier

test-unit: clean fmt build
	go test -v --race ./...

.PHONY: all build container clean fmt test-unit
