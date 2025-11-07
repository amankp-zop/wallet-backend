FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker

FROM alpine:latest

WORKDIR /

COPY --from=builder /api /api
COPY --from=builder /worker /worker

COPY ./configs/ /configs

EXPOSE 8080