# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.17-alpine as builder

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy
RUN go build


FROM alpine
# ENV TZ=Europe/Moscow
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /app/httpdumper ./app

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production
ENV GODEBUG=netdns=go+1

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 8397

# Uncomment to run the migrations before running the binary:
CMD /bin/app
#CMD exec /bin/app
