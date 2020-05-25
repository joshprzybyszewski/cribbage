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

# These happen to coincide with a local myql server with a root user and no password
ENV DSN_HOST=mysql
ENV DSN_USER=root
ENV DSN_PASSWORD=

CMD go run main.go \
    -restPort=8081 \
    -dsn_host=$DSN_HOST \
    -dsn_user=$DSN_USER \
    -dsn_password=$DSN_PASSWORD