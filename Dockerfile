FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /app/subscriptions-api ./cmd/subscriptions-api

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/subscriptions-api /app/subscriptions-api
COPY --from=build /app/configs /app/configs
COPY --from=build /app/migrations /app/migrations

RUN chmod +x /app/subscriptions-api

EXPOSE 8080

ENTRYPOINT ["/app/subscriptions-api"]