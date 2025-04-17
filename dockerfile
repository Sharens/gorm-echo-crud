FROM golang:1.24 AS builder


WORKDIR /build

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .



RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -tags osusergo,netgo -o /gorm-echo-crud .



FROM alpine:latest


WORKDIR /app

COPY --from=builder /gorm-echo-crud /app/gorm-echo-crud

EXPOSE 1323

CMD ["/app/gorm-echo-crud"]
