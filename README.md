# gorm-echo-crud
Prosty projekt API sklepu stworzony w języku Go z wykorzystaniem frameworka Echo oraz ORM-a GORM. Aplikacja umożliwia zarządzanie produktami, kategoriami oraz koszykiem. Przykład implementacji relacji w bazie danych oraz operacji CRUD na modelach.

## Wymagania projektu

* ❌ 3.0 Należy stworzyć aplikację we frameworki echo w j. Go, która będzie miała kontroler Produktów zgodny z CRUD
* ❌ 3.5 Należy stworzyć model Produktów wykorzystując gorm oraz wykorzystać model do obsługi produktów (CRUD) w kontrolerze (zamiast listy)
* ❌ 4.0 Należy dodać model Koszyka oraz dodać odpowiedni endpoint
* ❌ 4.5 Należy stworzyć model kategorii i dodać relację między kategorią, a produktem
* ❌ 5.0 pogrupować zapytania w gorm’owe scope'y

## Wymagania techniczne

Do uruchomienia aplikacji potrzebne są:

*   Zainstalowana aktualna wersja języka **Go** (np. 1.18 lub nowsza).
*   Pakiety Go zarządzane przez Go Modules (plik `go.mod`). Główne zależności to:
    *   `github.com/labstack/echo/v4` (framework Echo)
    *   `gorm.io/gorm` (GORM ORM)
    *   Odpowiedni sterownik bazy danych dla GORM (np. `gorm.io/driver/sqlite`, `gorm.io/driver/postgres`, w zależności od konfiguracji projektu).

Pakiety zostaną pobrane automatycznie przy pierwszym budowaniu lub uruchomieniu projektu za pomocą komend `go build` lub `go run`, lub ręcznie za pomocą `go mod tidy`.

## Uruchomienie aplikacji

1.  **(Opcjonalnie) Pobierz zależności:**
    Jeśli chcesz pobrać zależności przed pierwszym uruchomieniem:
    ```bash
    go mod tidy
    ```

2.  **Uruchom aplikację:**
    ```bash
    go run main.go
    ```

3.  **Dostęp do aplikacji:**
    Po pomyślnym uruchomieniu, serwer Echo nasłuchuje na określonym porcie pod adresem: http://localhost:1323
