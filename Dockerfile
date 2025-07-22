FROM golang:1.23-alpine

RUN apk add --no-cache postgresql-client

# WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main . 

CMD ["./main"]