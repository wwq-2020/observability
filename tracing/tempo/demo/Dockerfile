FROM golang:alpine

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app .

CMD ["./app"]