FROM golang:alpine AS builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update add tzdata ffmpeg

WORKDIR /app

EXPOSE 5000

COPY --from=builder /app/engine /app/config.json /app/

CMD mkdir files && \
    /app/engine