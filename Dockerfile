FROM golang:1.21 AS build

RUN mkdir /go/src/sai-scout
WORKDIR /go/src/sai-scout

COPY . .

RUN go build -o /sai-scout ./cmd/sai-scout/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

COPY --from=build /sai-scout /sai-scout

ENTRYPOINT ["/sai-scout"]