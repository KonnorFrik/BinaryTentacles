# === Build stage === #
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /go/src/
COPY go.mod go.sum /go/src/
RUN go mod download && apk add make
COPY . /go/src/
RUN make build name=order_service ver=1


# === Final stage === #
FROM golang:1.24-alpine3.21
WORKDIR /go/
EXPOSE 8888

COPY --from=builder /go/src/cmd/order_service/v1/order_service /go/bin/
RUN adduser -S user
USER user

CMD [ "bin/order_service" ]

