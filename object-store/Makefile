.PHONY: clean
clean: ## clean up previous builds
	@echo "\n cleaning previous builds"
	@rm -f bin/object-store

.PHONY: build
build: clean
	@echo  "\n Building object store binary"
	@go build -o bin/object-store ./cmd/...
	@echo "\n Build is successful built binary is at bin/object-store"

.PHONY: run
run:
	./bin/object-store

.PHONY: test
test:
	@go test ./dbstore/...