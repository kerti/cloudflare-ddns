MAJOR = 0
MINOR = 1

clean:
	@echo "Cleaning up project directory..." && \
	rm -rf .cover && \
	rm -f cloudflare-ddns && \
	rm -rf build

test:
	@echo "Running tests..." && \
	go test ./... -race

cover:
	@echo "Running tests and generating test coverage profile..." && \
	bash coverage.sh

cover-html:
	@echo "Running tests and displaying HTML test coverage profile..." && \
	bash coverage.sh --html

run-local:
	@echo "Building and running cloudflare-ddns locally..." && \
	go build -ldflags="-X main.majVersion=${MAJOR} -X main.minVersion=${MINOR} -X main.buildNum=$$(git rev-parse --short HEAD) -X main.verSuffix=local -s -w" && ./cloudflare-ddns --config config-local.yaml

run-dev:
	@echo "Building and running cloudflare-ddns dev version locally..." && \
	bash build-dev.sh --major=$(MAJOR) --minor=$(MINOR) && \
	./build/dev/cloudflare-ddns_darwin_amd64 --config config-local.yaml

run-rel:
	@echo "Building and running cloudflare-ddns release version locally..." && \
	bash build-release.sh --major=$(MAJOR) --minor=$(MINOR) && \
	./build/release/cloudflare-ddns_darwin_amd64 --config config-local.yaml
