FROM golang:1.25 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/scheduler .

FROM ubuntu:latest
WORKDIR /app

COPY --from=builder /out/scheduler /app/scheduler
COPY web /app/web

EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=

VOLUME ["/data"]
CMD ["/app/scheduler"]
