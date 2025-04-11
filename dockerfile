# Budowanie
FROM golang:1.24 AS builder

# Ustaw katalog roboczy w etapie budowania
WORKDIR /app

# Skopiuj pliki zależności
COPY go.mod go.sum ./

# Pobierz zależności
RUN go mod download

# Skopiuj resztę kodu źródłowego
COPY . .

# Zbuduj aplikację statycznie (zalecane dla alpine)
# Wynikowa binarka to /app/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/app .

# Tworzenie alpine
FROM alpine:latest

# Ustawianie katalogu roboczego
WORKDIR /app

COPY --from=builder /app/app .

# Wystaw port, na którym nasłuchuje aplikacja
EXPOSE 1323

# Uruchamia plik wykonywalny './app' znajdujący się w WORKDIR (/app)
CMD ["./app"]
