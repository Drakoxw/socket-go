FROM golang:1.20-alpine AS builder

WORKDIR /go/src

COPY ./go.mod ./go.sum ./

ENV GO111MODULE=on

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./app ./main.go

FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app /go/src/app

CMD ["/go/src/app"]

