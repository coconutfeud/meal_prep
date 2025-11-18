FROM golang:1.24 AS builder
ENV CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends build-essential gcc 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/

FROM debian:stable-slim
WORKDIR /root/
COPY --from=builder /app/app .
COPY --from=builder /app/public ./public
EXPOSE 80
CMD ["./app"]