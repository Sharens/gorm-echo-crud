FROM golang:1.24 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/run .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/run .
EXPOSE 1323

CMD ["./run"]
