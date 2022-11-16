## Build

FROM golang:1.17-alpine AS build
ENV GOPATH /go
ENV GOBIN $GOPATH/bin

COPY . /go/src/github.com/doka-guide/api/
WORKDIR /go/src/github.com/doka-guide/api/
RUN go mod download
RUN go build -o /app/api
COPY .env /app/

# Deploy

FROM alpine:latest
ARG APP_PORT

VOLUME /app
WORKDIR /app
COPY --from=build /app ./
EXPOSE $APP_PORT
ENTRYPOINT ["/app/api"]