FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o api cmd/api/main.go

FROM alpine:latest
COPY --from=builder /app/api /api
EXPOSE 8080
CMD ["/api"]
