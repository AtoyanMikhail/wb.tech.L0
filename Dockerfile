FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/static ./static

COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8081

CMD ["./main"] 