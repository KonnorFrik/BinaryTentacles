# === Build stage === #
FROM golang:1.24-alpine3.21 AS builder
COPY go.mod go.sum /go/src/
WORKDIR /go/src/
RUN go mod download && apk add make
COPY . /go/src/
RUN make build name=spot_instrument ver=1


# === Final stage === #
FROM golang:1.24-alpine3.21
WORKDIR /go/
EXPOSE 9999

COPY --from=builder /go/src/cmd/spot_instrument/v1/spot_instrument /go/bin/
RUN adduser -S user
USER user

CMD [ "bin/spot_instrument" ]

