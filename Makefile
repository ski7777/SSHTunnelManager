.PHONY: deps
deps:
	go mod vendor

.PHONY: build
build: output/SSHTunnelManager

.PHONY: build-ar71xx
build-ar71xx: output/SSHTunnelManager.ar71xx

.PHONY: build-all
build-all: build build-ar71xx 

.PHONY: prepare-build
prepare-build:
	mkdir -p output

.PHONY: clean
clean:
	rm -rf output

output/SSHTunnelManager: prepare-build
	go build -trimpath -ldflags="-s -w" -o output/SSHTunnelManager cmd/main.go

output/SSHTunnelManager.ar71xx: prepare-build
	GOARCH=mips GOMIPS=softfloat CC=mips-linux-gnu-gcc go build -trimpath -ldflags="-s -w" -o output/SSHTunnelManager.ar71xx cmd/main.go

