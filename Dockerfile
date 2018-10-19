FROM golang:1.11 as builder
WORKDIR /go/src/github.com/filiphe/goba/
COPY main.go .
COPY main_test.go .
COPY go.mod .
COPY go.sum .
COPY drinks.json .
RUN GO111MODULE=on go test
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -a -installsuffix cgo -o goba .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/filiphe/goba/goba .
COPY drinks.json .
EXPOSE 3333
CMD ["./goba"]

