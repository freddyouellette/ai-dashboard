.PHONY: *

COVERAGE_REQUIREMENT := 90

clear-build:
	rm -rf ./bin

build-frontend:
	cd ./web && npm i && npm run build && cd ../

build-backend:
	mkdir -p ./bin && CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -v -o ./bin ./cmd/api/main.go

build-plugins:
	@for file in $$(find ./plugins -name '*.go'); do \
		if grep -q "^package main" "$$file"; then \
			echo "Building plugin: $$file"; \
			go build -buildmode=plugin -modfile go.mod -o "$${file%.go}.so" "$$file"; \
		fi; \
	done

# build plugins to work with delve debugger
build-plugins-debug:
	@for file in $$(find ./plugins -name '*.go'); do \
		if grep -q "^package main" "$$file"; then \
			echo "Building plugin: $$file"; \
			go build -gcflags=all="-N -l" -buildmode=plugin -modfile go.mod -o "$${file%.go}.so" "$$file"; \
		fi; \
	done

build: clear-build build-plugins build-backend build-frontend

test:
	go test ./... -coverprofile=coverage.out

format:
	go fmt ./...

lint:
	golangci-lint run

ci-lint:
	actionlint

coverage:
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@COVERAGE=`go tool cover -func=coverage.out | grep "^total:" | grep -Eom 1 '[0-9]+' | head -1`;\
	echo Coverage report available at ./coverage.html;\
	if [ "$$COVERAGE" -lt "${COVERAGE_REQUIREMENT}" ]; then\
		echo "Test coverage $$COVERAGE% does not meet minimum ${COVERAGE_REQUIREMENT}% requirement";\
		exit 1;\
	else\
		echo "Test Coverage $$COVERAGE% (OK)";\
	fi

vulns:
	govulncheck ./...

pipeline: build test lint ci-lint coverages

app: build
	./bin/main

backend-dev: 
	FRONTEND=false nodemon --watch './**/*.go' --signal SIGTERM --exec 'make build-plugins && CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go run ./cmd/api/main.go'

frontend-dev:
	cd ./web && npm i && npm run dev && cd ../

app-dev: build-plugins
	make -j 2 backend-dev frontend-dev

app-local:
	DEVCONTAINER_COMMAND="make app" docker-compose up

# for custom make jobs
-include Makefile.local.mk