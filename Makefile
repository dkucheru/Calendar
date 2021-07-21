test_coverage:
	go test ./service/ -coverprofile=coverage.out

html_test:
	go tool cover -html=coverage.out

run:
	go run ./cmd/