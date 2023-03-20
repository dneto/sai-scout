FROM golang:1.20-alpine AS build

RUN mkdir /go/src/sai-scout
WORKDIR /go/src/sai-scout

COPY . .

RUN go build -o /sai-scout ./cmd/sai-scout/main.go

FROM alpine:latest

COPY --from=build /go/src/sai-scout/cards.json /cards.json
COPY --from=build /sai-scout /sai-scout

ENTRYPOINT ["/sai-scout"]