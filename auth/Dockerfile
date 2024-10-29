FROM golang:1.23.2 AS builder

WORKDIR /build

COPY . .

RUN ls

RUN cd auth && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /auth

FROM golang:1.23.2

WORKDIR /app

COPY --from=builder /auth auth
COPY --from=builder /build/auth/migrations migrations
COPY --from=builder /build/auth/templates templates

CMD ["/app/auth"]