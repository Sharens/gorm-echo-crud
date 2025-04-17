# gorm-echo-crud
Prosty projekt API sklepu stworzony w języku Go z wykorzystaniem frameworka Echo oraz ORM-a GORM. Aplikacja umożliwia zarządzanie produktami, kategoriami oraz koszykiem. Przykład implementacji relacji w bazie danych oraz operacji CRUD na modelach.

## Wymagania projektu

* ✅ 3.0 Należy stworzyć aplikację we frameworki echo w j. Go, która będzie miała kontroler Produktów zgodny z CRUD
* ✅ 3.5 Należy stworzyć model Produktów wykorzystując gorm oraz wykorzystać model do obsługi produktów (CRUD) w kontrolerze (zamiast listy)
* ✅ 4.0 Należy dodać model Koszyka oraz dodać odpowiedni endpoint
* ✅ 4.5 Należy stworzyć model kategorii i dodać relację między kategorią, a produktem
* ✅ 5.0 pogrupować zapytania w gorm’owe scope'y

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

## Budowanie obrazu Docker

1.  **Budowanie obrazu Dockerowego**
    ```
    docker build -t gorm_echo_crud .
    ```

2.  **Odpalanie aplikacji z obrazu na port 1323**
    ```
    docker run -p 1323:1323 gorm_echo_crud:latest 
    ```


## Testowanie API (curl)

Poniższe komendy `curl` pozwalają przetestować wszystkie endpointy API. Zakładają, że serwer działa na `http://localhost:1323`.

**Ważne:** ID tworzonych zasobów (kategorie, produkty, koszyki, elementy koszyka) są nadawane przez bazę danych. W poniższych przykładach zakładamy pewne ID (np. 1, 2), ale w rzeczywistości mogą być inne. Należy używać ID zwróconych przez API w poprzednich krokach.

### Kategorie (`/categories`)

1.  **Utwórz nową kategorię:**
    ```bash
    curl -X POST http://localhost:1323/categories -H "Content-Type: application/json" -d '{"name": "Elektronika"}'
    # Oczekiwana odpowiedź: JSON z utworzoną kategorią (zanotuj ID, np. 1)
    ```
    ```bash
    curl -X POST http://localhost:1323/categories -H "Content-Type: application/json" -d '{"name": "Książki"}'
    # Oczekiwana odpowiedź: JSON z utworzoną kategorią (zanotuj ID, np. 2)
    ```
    ```bash
    # Spróbuj utworzyć kategorię o tej samej nazwie
    curl -i -X POST http://localhost:1323/categories -H "Content-Type: application/json" -d '{"name": "Elektronika"}'
    # Oczekiwana odpowiedź: Status 409 Conflict
    ```

2.  **Pobierz wszystkie kategorie:**
    ```bash
    curl http://localhost:1323/categories
    # Oczekiwana odpowiedź: Lista kategorii w JSON
    ```

3.  **Pobierz kategorię po ID (np. ID=1):**
    ```bash
    curl http://localhost:1323/categories/1
    # Oczekiwana odpowiedź: Dane kategorii o ID 1 w JSON
    ```

4.  **Zaktualizuj kategorię (np. ID=1):**
    ```bash
    curl -X PUT http://localhost:1323/categories/1 -H "Content-Type: application/json" -d '{"name": "AGD i Elektronika"}'
    # Oczekiwana odpowiedź: Zaktualizowane dane kategorii
    ```

5.  **Usuń kategorię (np. ID=2):**
    *   *Uwaga: Usunięcie kategorii powiedzie się tylko, jeśli nie ma w niej żadnych produktów.*
    ```bash
    # Najpierw upewnij się, że kategoria jest pusta
    curl -X DELETE http://localhost:1323/categories/2
    # Oczekiwana odpowiedź: Brak zawartości (Status 204)
    ```
    ```bash
    # Spróbuj usunąć kategorię, która zawiera produkty (np. ID=1, jeśli dodano do niej produkt)
    curl -i -X DELETE http://localhost:1323/categories/1
    # Oczekiwana odpowiedź: Status 409 Conflict (jeśli zawiera produkty)
    ```

### Produkty (`/products`)

*   *Wymaga istnienia kategorii.*

1.  **Utwórz nowy produkt (w kategorii o ID=1):**
    ```bash
    curl -X POST http://localhost:1323/products -H "Content-Type: application/json" -d '{"name": "Laptop", "price": 3500.99, "category_id": 1}'
    # Oczekiwana odpowiedź: JSON z utworzonym produktem (zanotuj ID, np. 1)
    ```
    ```bash
    curl -X POST http://localhost:1323/products -H "Content-Type: application/json" -d '{"name": "Myszka", "price": 79.00, "category_id": 1}'
    # Oczekiwana odpowiedź: JSON z utworzonym produktem (zanotuj ID, np. 2)
    ```
    ```bash
    # Spróbuj utworzyć produkt z nieistniejącą kategorią
    curl -i -X POST http://localhost:1323/products -H "Content-Type: application/json" -d '{"name": "Duch", "price": 10, "category_id": 999}'
    # Oczekiwana odpowiedź: Status 400 Bad Request "Invalid input: Category not found"
    ```

2.  **Pobierz wszystkie produkty:**
    ```bash
    curl http://localhost:1323/products
    # Oczekiwana odpowiedź: Lista produktów w JSON (z danymi kategorii)
    ```

3.  **Pobierz produkt po ID (np. ID=1):**
    ```bash
    curl http://localhost:1323/products/1
    # Oczekiwana odpowiedź: Dane produktu o ID 1 w JSON (z danymi kategorii)
    ```

4.  **Zaktualizuj produkt (np. ID=1, zmień cenę i kategorię na ID=2, jeśli istnieje):**
    ```bash
    # Upewnij się, że kategoria o ID 2 istnieje
    curl -X PUT http://localhost:1323/products/1 -H "Content-Type: application/json" -d '{"name": "Laptop Gamingowy", "price": 4200, "category_id": 2}'
    # Oczekiwana odpowiedź: Zaktualizowane dane produktu
    ```

5.  **Usuń produkt (np. ID=2):**
    ```bash
    curl -X DELETE http://localhost:1323/products/2
    # Oczekiwana odpowiedź: Brak zawartości (Status 204)
    ```

### Koszyk (`/cart`)

*   *Wymaga istnienia produktów.*

1.  **Utwórz nowy koszyk:**
    ```bash
    curl -X POST http://localhost:1323/cart
    # Oczekiwana odpowiedź: JSON z utworzonym koszykiem (zanotuj ID koszyka, np. 1)
    ```

2.  **Dodaj produkt do koszyka (np. produkt ID=1 do koszyka ID=1):**
    ```bash
    # Upewnij się, że produkt o ID 1 istnieje
    curl -X POST http://localhost:1323/cart/1/items -H "Content-Type: application/json" -d '{"product_id": 1, "quantity": 2}'
    # Oczekiwana odpowiedź: JSON z dodanym elementem koszyka (zanotuj item_id, np. 1)
    ```
    ```bash
    # Dodaj ten sam produkt ponownie (zwiększy ilość)
    curl -X POST http://localhost:1323/cart/1/items -H "Content-Type: application/json" -d '{"product_id": 1, "quantity": 1}'
    # Oczekiwana odpowiedź: JSON ze zaktualizowanym elementem koszyka (quantity powinno być 3)
    ```
    ```bash
    # Spróbuj dodać nieistniejący produkt
    curl -i -X POST http://localhost:1323/cart/1/items -H "Content-Type: application/json" -d '{"product_id": 998, "quantity": 1}'
    # Oczekiwana odpowiedź: Status 404 Not Found "Product not found"
    ```

3.  **Pobierz zawartość koszyka (np. ID=1):**
    ```bash
    curl http://localhost:1323/cart/1
    # Oczekiwana odpowiedź: Dane koszyka z listą elementów (wraz z danymi produktów)
    ```

4.  **Usuń element z koszyka (np. element o ID=1 z koszyka ID=1):**
    ```bash
    curl -X DELETE http://localhost:1323/cart/1/items/1
    # Oczekiwana odpowiedź: Brak zawartości (Status 204)
    ```

5.  **Usuń cały koszyk (np. ID=1):**
    ```bash
    curl -X DELETE http://localhost:1323/cart/1
    # Oczekiwana odpowiedź: Brak zawartości (Status 204)
    ```
    ```bash
    # Spróbuj pobrać usunięty koszyk
    curl -i http://localhost:1323/cart/1
    # Oczekiwana odpowiedź: Status 404 Not Found
    ```

## Interfejs Webowy

Prosty interfejs do testowania endpointów produktów jest dostępny pod adresem:
`http://localhost:1323/test-ui`
*(Uwaga: Ten interfejs może nie obsługiwać jeszcze wszystkich nowych endpointów kategorii i koszyka).*

Strona główna z nawigacją:
`http://localhost:1323/`

