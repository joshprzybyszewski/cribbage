FROM golang:1.14.3-alpine3.11

ENV GOATH=/go
RUN mkdir -p $GOATH/src/github.com/joshprzybyszewski/cribbage
WORKDIR $GOATH/src/github.com/joshprzybyszewski/cribbage

COPY model model
COPY network network
COPY logic logic
COPY utils utils
COPY jsonutils jsonutils
COPY server server
COPY vendor vendor
COPY main.go main.go

EXPOSE 80

ARG dsn_host
ARG dsn_user
ARG dsn_password

CMD go run main.go \
    -restPort=8081 \
    -dsn_host=redacted \
    -dsn_user=redacted \
    -dsn_password=redacted