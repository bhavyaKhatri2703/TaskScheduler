FROM golang:1.24.4 as builder

WORKDIR /app

# Set Go proxy to direct (no proxy)
ENV GOPROXY="https://goproxy.io"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
