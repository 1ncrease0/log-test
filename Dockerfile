FROM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /log-parser ./cmd/log-parser

FROM alpine:3.20


COPY --from=build /log-parser /log-parser

WORKDIR /app

ENTRYPOINT ["/log-parser"]
