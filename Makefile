build:
	go build ./...
run: 
	go run ./cmd/main.go
setup:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest


lint:
	golangci-lint run ./... --timeout 10m

lint-fix:
	golangci-lint run ./... --fix

mocks:
	go generate ./...

test:
	ginkgo ./...

e2e-test:
	godotenv -f .env.staging ginkgo ./cmd

test-report:
	ginkgo --junit-report=report.xml  ./... 

it: export TEST_TYPE = IT
it:
	ginkgo --race ./internal/repository/...
