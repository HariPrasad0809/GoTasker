FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gotasker main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /gotasker .
COPY .env ./

EXPOSE 8080

CMD ["/app/gotasker"]