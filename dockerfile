FROM golang:1.24 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /gorm-echo-crud .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/view ./view
COPY --from=builder /gorm-echo-crud /app/gorm-echo-crud

EXPOSE 1323

CMD ["/app/gorm-echo-crud"]
