FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN go build -o /cmd/main.go



FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/ /app/cmd/

COPY --from=builder /app/web/ /app/web/

ENV TODO_PORT=:7540 TODO_DBFILE=../scheduler.db TODO_PASSWORD=12345

EXPOSE 7540

CMD ["/app/cmd/"]