ALL_PACKAGES=$(shell (go list ./... ))
OUT_DIR := ./out
BIN := ${OUT_DIR}/bin/zearch

.PHONY: build test race cover clean run lint help

build:
	rm -rf $(BIN)
	go build -o $(BIN) ./cmd/zearch

test:
	go test $(ALL_PACKAGES) -timeout 10s

race:
	go test $(ALL_PACKAGES) -timeout 10s -race

cover:
	rm -rf out/cover/*
	mkdir -p out/cover
	go test -short -cover -coverpkg=./... -coverprofile=out/cover/test.coverage ./...
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

lint:
	@echo "Running linter in docker container"
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run

docker:
	docker build -t jaimem88/zearch .
	docker run -ti --rm jaimem88/zearch

# Obtained from https://stackoverflow.com/questions/4219255/how-do-you-get-the-list-of-targets-in-a-makefile
help:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
