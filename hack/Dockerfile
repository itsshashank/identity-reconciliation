FROM golang:1.20

WORKDIR /app

COPY . ../.

RUN go mod download && go mod verify

RUN go build -o app /cmd/main.go

CMD ["./app"]
