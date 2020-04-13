MAJOR = 0
MINOR = 1

.PHONY: clean test cover cover-html build-local build-dev build-rel run-local run-dev run-rel

clean:
	@echo "Cleaning up project directory..." && \
	rm -rf .cover && \
	rm -f cloudflare-ddns && \
	rm -rf build

test: clean
	@echo "Running tests..." && \
	go test ./... -race

cover:
	@echo "Running tests and generating test coverage profile..." && \
	bash coverage.sh

cover-html:
	@echo "Running tests and displaying HTML test coverage profile..." && \
	bash coverage.sh --html

build-local: test
	@echo "Building cloudflare-ddns locally..." && \
	go build -ldflags="-X main.majVersion=${MAJOR} -X main.minVersion=${MINOR} -X main.buildNum=$$(git rev-parse --short HEAD) -X main.verSuffix=local -s -w"

build-dev: test
	@echo "Building cloudflare-ddns dev version locally..." && \
	bash build-dev.sh --major=$(MAJOR) --minor=$(MINOR)

build-rel: test
	@echo "Building cloudflare-ddns release version locally..." && \
	bash build-release.sh --major=$(MAJOR) --minor=$(MINOR)

run-local: build-local
	@echo "Running cloudflare-ddns locally..."
	./cloudflare-ddns --config config-local.yaml

run-dev: build-dev
	@echo "Running cloudflare-ddns dev version locally..." && \
	./build/dev/cloudflare-ddns_darwin_amd64 --config config-local.yaml

run-rel: build-rel
	@echo "Running cloudflare-ddns release version locally..." && \
	./build/release/cloudflare-ddns_darwin_amd64 --config config-local.yaml
