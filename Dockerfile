FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations
CMD ["./app"]