test_coverage:
	go test ./service/ -coverprofile=coverage.out

html_test:
	go tool cover -html=coverage.out

run:
	go run ./cmd/

build:
	CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

db-up:
	sudo docker run -dp 5432:5432 \
    --name db-container2 \
    -e POSTGRES_PASSWORD=xxxxx \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -e POSTGRES_USER=xxxx \
    -e POSTGRES_DB=calendar \
    -v /custom/mount:/var/lib/postgresql/data \
    postgres
