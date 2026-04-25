FROM golang:1.26.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api    cmd/api/main.go
RUN go build -o worker cmd/worker/main.go

# ---- API production image ----
FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/api /app/api
EXPOSE ${PORT}
CMD ["./api"]

# ---- Worker production image ----
FROM alpine:3.20.1 AS worker
WORKDIR /app
COPY --from=build /app/worker /app/worker
CMD ["./worker"]
