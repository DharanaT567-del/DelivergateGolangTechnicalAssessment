package tests

import (
	"net/http"
	"strconv"
	"testing"
)

func TestCreateAndDeleteCustomer(t *testing.T) {
	app := newTestApp(t)
	token := loginToken(t, app, "owner", "owner123")

	created := request(t, app, http.MethodPost, "/api/v1/customers",
		`{"first_name":"Jamie","last_name":"Lee","phone":"555-1234","email":"jamie@example.com","address":"Main Street"}`, token)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("create customer status = %d, want %d", created.StatusCode, http.StatusCreated)
	}
	body := responseBody(t, created)
	data := body["data"].(map[string]any)
	id := strconv.Itoa(int(data["id"].(float64)))

	invalid := request(t, app, http.MethodPost, "/api/v1/customers",
		`{"first_name":"Jamie","phone":"555-1234","email":"invalid-email"}`, token)
	if invalid.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid customer status = %d, want %d", invalid.StatusCode, http.StatusBadRequest)
	}
	invalid.Body.Close()

	deleted := request(t, app, http.MethodDelete, "/api/v1/customers/"+id, "", token)
	if deleted.StatusCode != http.StatusOK {
		t.Fatalf("delete customer status = %d, want %d", deleted.StatusCode, http.StatusOK)
	}
	deleted.Body.Close()

	fetched := request(t, app, http.MethodGet, "/api/v1/customers/"+id, "", token)
	if fetched.StatusCode != http.StatusNotFound {
		t.Fatalf("get soft-deleted customer status = %d, want %d", fetched.StatusCode, http.StatusNotFound)
	}
	fetched.Body.Close()
}
