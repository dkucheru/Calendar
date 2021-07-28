FROM golang:latest as builder
RUN mkdir /appCalendar
COPY go.mod go.sum /appCalendar/
WORKDIR /appCalendar
RUN go mod download
COPY . .

RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /appCalendar/main /usr/bin/calendar
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/calendar"]

