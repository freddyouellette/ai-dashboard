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

api:
	go run ./cmd/api/main.go

frontend:
	cd ./web && npm i && npm start && cd ../

app: build
	./bin/api

app-dev-api:
	FRONTEND=false nodemon --watch './**/*.go' --signal SIGTERM --exec go run ./cmd/api/main.go

app-dev-frontend:
	cd ./web && npm i && npm start && cd ../

app-dev:
	make -j 2 app-dev-api app-dev-frontend

app-local:
	DEVCONTAINER_COMMAND="make app" docker-compose up