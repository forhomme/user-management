FROM golang:1.20-alpine as builder
LABEL stage=builder

RUN apk add --no-cache make git build-base

ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
ENV GOPRIVATE=github.com/forhomme/*
ENV CGO_ENABLED=0

COPY .  /usr/app/go
WORKDIR /usr/app/go
RUN go env
RUN go build -ldflags="-s -w" -o users-management

FROM ubuntu:20.04
WORKDIR /app

COPY --from=builder /usr/app/go/user-management /app/
COPY --from=builder /usr/app/go/config.yaml /app/
COPY --from=builder /usr/app/go/migration/ /app/migration

RUN apt-get update && apt-get upgrade -y && apt-get install tzdata ca-certificates -y
ENV TZ="Asia/Singapore"
EXPOSE 8081

CMD ["./user-management"]
