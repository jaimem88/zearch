ALL_PACKAGES=$(shell (go list ./...))
OUT_DIR := ./out
BIN := ${OUT_DIR}/bin/zearch

.PHONY: build test cover clean run

build:
	rm -rf $(BIN)
	go build -o $(BIN) ./cmd/zearch

test:
	go test $(ALL_PACKAGES) -timeout 10s

cover:
	rm -rf out/cover/*
	mkdir -p out/cover
	go test -short -cover -coverpkg=$(ALL_PACKAGES) -coverprofile=out/cover/test.coverage $(ALL_PACKAGES)
	go tool cover -html out/cover/test.coverage -o out/cover/coverage.html
	@echo ""
	@echo "=====> Total test coverage: <====="
	@echo ""
	go tool cover -func out/cover/test.coverage

clean:
	echo "Removing out/"
	rm -rf out/*

run: build
	$(BIN) ${ARGS}
