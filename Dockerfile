# Build stage
FROM golang:1.24 as builder

WORKDIR /
COPY . .

RUN go mod download

RUN go build -o app .

# Final image
FROM debian:bullseye-slim

WORKDIR /
COPY --from=builder /app .

EXPOSE 1323
ENTRYPOINT ["./app"]
