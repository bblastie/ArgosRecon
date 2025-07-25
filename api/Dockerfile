FROM golang:1.23-alpine

RUN apk add --no-cache postgresql-client
WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main . 

CMD ./main