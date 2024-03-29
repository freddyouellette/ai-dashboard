.PHONY: *

COVERAGE_REQUIREMENT := 90

build-frontend:
	cd ./web && npm i && npm run build && cd ../

build-backend:
	mkdir -p ./bin && go build -v -o ./bin ./...

build: build-frontend build-backend

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
	./bin/api

backend-dev:
	FRONTEND=false nodemon --watch './**/*.go' --signal SIGTERM --exec go run ./cmd/api/main.go

frontend-dev:
	cd ./web && npm i && npm run dev && cd ../

app-dev:
	make -j 2 backend-dev frontend-dev

app-local:
	DEVCONTAINER_COMMAND="make app" docker-compose up