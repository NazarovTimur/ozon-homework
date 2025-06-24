.PHONY: run-all test-coverage cyclo cognit


run-all:
	go run ./cmd/server/main.go

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

cyclo:
	gocyclo -over 10 .

cognit:
	gocognit -over 10 .