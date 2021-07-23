test_coverage:
	go test ./service/ -coverprofile=coverage.out

html_test:
	go tool cover -html=coverage.out

run:
	go run ./cmd/

docker-build:
	cd /appCalendar/cmd/ && CGO_ENABLED=0 GOOS=linux go build