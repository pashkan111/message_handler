FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go clean -modcache
RUN go mod tidy
RUN go build -o main .

CMD ["./main"]