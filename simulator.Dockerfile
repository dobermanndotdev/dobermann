FROM golang:1.21 as builder

ARG VERSION

WORKDIR /usr/app
COPY . ./

ENV CGO_ENABLED=0
RUN go build -buildvcs=false -o bin/service -ldflags="-X main.Version=${VERSION}" ./cmd/endpoint_simulator

FROM alpine
WORKDIR /usr/app
COPY --from=builder /usr/app/bin/service ./service

ENTRYPOINT ["./service"]