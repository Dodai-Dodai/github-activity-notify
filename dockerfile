# build stage
FROM golang:latest AS builder

WORKDIR /work

COPY main.go returnHour.go sendLine.go go.mod go.sum ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# exec stage
FROM gcr.io/distroless/static-debian11:latest 

WORKDIR /root/

COPY --from=builder /work/main .

ENTRYPOINT ["./main"]
