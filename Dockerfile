FROM golang:1.25.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/bin/server /app/server

EXPOSE 8080

CMD ["/app/server"]
