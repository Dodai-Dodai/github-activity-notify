FROM golang:latest

WORKDIR /work

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go build -o main main.go

ENTRYPOINT ["./main"]