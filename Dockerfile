FROM golang:1.13.1

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./cmd ./cmd

RUN go build -o titsbot ./cmd

CMD ["./titsbot"]