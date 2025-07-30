FROM golang:1.24.5-alpine

WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subdomainParser . 

CMD ./subdomainParser