FROM golang:1.25 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/scheduler .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /out/scheduler /app/scheduler
COPY web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=12345

VOLUME ["/data"]
CMD ["/app/scheduler"]
