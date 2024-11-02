# syntax=docker/dockerfile:1

FROM node:20 AS tailwind

WORKDIR /app

COPY web/package.json ./web/
COPY web/package-lock.json ./web/

WORKDIR /app/web

RUN npm install

COPY web ./

RUN npm run build

FROM golang:1.22

ARG HTTP_PORT=8080

WORKDIR /app

COPY --from=tailwind /app/web/static ./web/static

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener/main.go

EXPOSE ${HTTP_PORT}

CMD ["./url-shortener"]