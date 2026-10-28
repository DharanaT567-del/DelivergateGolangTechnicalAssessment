package tests

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestCreateUpdateDeleteItem(t *testing.T) {
	app := newTestApp(t)
	token := loginToken(t, app, "owner", "owner123")

	created := request(t, app, http.MethodPost, "/api/v1/items",
		`{"name":"Chocolate","category":"Ice Cream","quantity":10,"price":4.25,"description":"Dark chocolate"}`, token)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("create item status = %d, want %d", created.StatusCode, http.StatusCreated)
	}
	createdBody := responseBody(t, created)
	data, ok := createdBody["data"].(map[string]any)
	if !ok {
		t.Fatalf("created response has no data object: %#v", createdBody)
	}
	id := strconv.Itoa(int(data["id"].(float64)))

	updated := request(t, app, http.MethodPut, "/api/v1/items/"+id,
		`{"name":"Chocolate","category":"Frozen Dessert","quantity":8,"price":4.75}`, token)
	if updated.StatusCode != http.StatusOK {
		t.Fatalf("update item status = %d, want %d", updated.StatusCode, http.StatusOK)
	}
	updated.Body.Close()

	deleted := request(t, app, http.MethodDelete, "/api/v1/items/"+id, "", token)
	if deleted.StatusCode != http.StatusOK {
		t.Fatalf("delete item status = %d, want %d", deleted.StatusCode, http.StatusOK)
	}
	deleted.Body.Close()

	fetched := request(t, app, http.MethodGet, "/api/v1/items/"+id, "", token)
	if fetched.StatusCode != http.StatusNotFound {
		t.Fatalf("get soft-deleted item status = %d, want %d", fetched.StatusCode, http.StatusNotFound)
	}
	fetched.Body.Close()
}

func TestSearchItemsAndValidation(t *testing.T) {
	app := newTestApp(t)
	token := loginToken(t, app, "owner", "owner123")
	for _, name := range []string{"Chocolate", "Vanilla"} {
		response := request(t, app, http.MethodPost, "/api/v1/items",
			`{"name":"`+name+`","category":"Ice Cream","quantity":3,"price":2.5}`, token)
		if response.StatusCode != http.StatusCreated {
			t.Fatalf("create item %s status = %d", name, response.StatusCode)
		}
		response.Body.Close()
	}

	search := request(t, app, http.MethodGet, "/api/v1/items?search=Choc&page=1&limit=1&sort=-name", "", token)
	if search.StatusCode != http.StatusOK {
		t.Fatalf("search status = %d, want %d", search.StatusCode, http.StatusOK)
	}
	body := responseBody(t, search)
	data := body["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) != 1 || !strings.Contains(items[0].(map[string]any)["name"].(string), "Chocolate") {
		t.Fatalf("unexpected search response items: %#v", items)
	}
	if data["total"].(float64) != 1 {
		t.Fatalf("search total = %v, want 1", data["total"])
	}

	invalid := request(t, app, http.MethodPost, "/api/v1/items",
		`{"name":" ","category":"Ice Cream","quantity":-1,"price":0}`, token)
	if invalid.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid item status = %d, want %d", invalid.StatusCode, http.StatusBadRequest)
	}
	invalid.Body.Close()
}
