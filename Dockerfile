FROM golang:1.25.3-alpine AS build

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/bin/server /app/server
COPY web/templates /app/web/templates

EXPOSE 8080

CMD ["/app/server"]
