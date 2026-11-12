.PHONY: build test lint proto chaos cluster

build:
	go build ./...

test:
	go test -race ./...

lint:
	gofmt -l . | (! grep .)
	go vet ./...
	staticcheck ./...

proto:
	./scripts/gen-proto.sh

# ponytail: stubs until phase 6 / phase 4 land
chaos:
	@echo "chaos: not implemented until phase 6"

cluster:
	@echo "cluster: not implemented until phase 4"
