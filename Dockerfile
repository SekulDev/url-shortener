# syntax=docker/dockerfile:1

FROM node:20

WORKDIR /app

COPY web/package.json ./web/
COPY web/package-lock.json ./web/

RUN cd web && npm install
RUN cd web && npm run build

FROM golang:1.22

ARG HTTP_PORT=8080

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener/main.go

EXPOSE ${HTTP_PORT}

CMD ["./url-shortener"]