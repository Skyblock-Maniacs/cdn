# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cdn ./cmd/cdn

ENTRYPOINT [ "/cdn" ]
