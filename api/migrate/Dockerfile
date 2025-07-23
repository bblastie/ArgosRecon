FROM golang:1.23-alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client
COPY --from=builder /go/bin/goose /usr/local/bin/goose
WORKDIR /migrations
ENTRYPOINT ["sh", "-c", "goose postgres \"$DATABASE_URL\" up"]