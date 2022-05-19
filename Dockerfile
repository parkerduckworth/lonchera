# syntax=docker/dockerfile:1
FROM golang:1.17 as builder
COPY . /build
WORKDIR /build
RUN  go mod tidy -go=1.16 && go mod tidy -go=1.17 && \
     gofmt -w . && \
     go build -o service .

FROM alpine:latest as release
RUN apk --update --no-cache add ca-certificates gcompat
RUN addgroup -S servicegroup && adduser -S serviceuser -G servicegroup

COPY --from=builder /build/env /build/env
COPY --from=builder /build/service /build

USER serviceuser
WORKDIR /build

FROM release as dev
ENV GO_ENV=dev
CMD ["./service"]

FROM release as staging
ENV GO_ENV=staging
CMD ["./service"]

FROM release as prod
ENV GO_ENV=prod
CMD ["./service"]
