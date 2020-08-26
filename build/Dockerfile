FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG project
ENV PROJECT $project

ENV GOOS linux
ENV CGO_ENABLED 0
ENV GOARCH amd64
RUN go build -o ./out/app $project


FROM alpine:3.9 
COPY --from=build_base /tmp/app/out/app /app/app
EXPOSE 8080
CMD ["/app/app"]