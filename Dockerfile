#FROM golang:1.12 AS build
#COPY . /fgservice
#WORKDIR /fgservice
#RUN apt update && apt install ca-certificates libgnutls30 -y
#RUN go get -d
# or FROM golang:alpine or some other base depending on need
#FROM alpine:latest AS runtime
#this seems dumb, but the libc from the build stage is not the same as the alpine libc
#create a symlink to where it expects it since they are compatable. https://stackoverflow.com/a/35613430/3105368
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#WORKDIR /fgservice
#COPY --from=build /fgservice/menu-frontend-go .
# Build image
########################
#FROM golang:1.21.4-alpine3.18 as builder

#WORKDIR /var/tmp/app

#RUN apk add git

# copy artifacts into the container
#ADD ./cmd ./cmd
#ADD ./go.mod ./go.mod
#ADD ./go.sum ./go.sum
#ADD ./internal ./internal

# Build the app
#RUN go build -o .build/app ./cmd/cli

# Final image
########################
#FROM alpine:3.18.4
#
#WORKDIR /opt/app
#
#COPY --from=builder /var/tmp/app/.build/app .

#CMD [ "./app" ]
# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main .
EXPOSE 8080
CMD ["/app/main"]