FROM golang:alpine as builder

WORKDIR /go/src/app

COPY . .

RUN go get -v
RUN go install -v

RUN go build main.go




FROM alpine

COPY --from=builder /go/src/app /app
WORKDIR /app
CMD ["./main"]