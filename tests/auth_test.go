package tests

import (
	"net/http"
	"testing"
	"time"

	"icecream-shop/utils"
)

func TestAuthentication(t *testing.T) {
	app := newTestApp(t)

	invalid := request(t, app, http.MethodPost, "/api/v1/auth/login",
		`{"username":"owner","password":"wrong"}`, "")
	if invalid.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid login status = %d, want %d", invalid.StatusCode, http.StatusUnauthorized)
	}
	invalid.Body.Close()

	token := loginToken(t, app, "owner", "owner123")
	protected := request(t, app, http.MethodGet, "/api/v1/items", "", token)
	if protected.StatusCode != http.StatusOK {
		t.Fatalf("authenticated request status = %d, want %d", protected.StatusCode, http.StatusOK)
	}
	protected.Body.Close()

	missing := request(t, app, http.MethodGet, "/api/v1/items", "", "")
	if missing.StatusCode != http.StatusUnauthorized {
		t.Fatalf("missing token status = %d, want %d", missing.StatusCode, http.StatusUnauthorized)
	}
	missing.Body.Close()
}

func TestRoleAuthorization(t *testing.T) {
	app := newTestApp(t)
	managerToken := loginToken(t, app, "manager", "manager123")
	cashierToken := loginToken(t, app, "cashier", "cashier123")
	ownerToken := loginToken(t, app, "owner", "owner123")

	item := `{"name":"Vanilla","category":"Ice Cream","quantity":5,"price":3.5}`
	created := request(t, app, http.MethodPost, "/api/v1/items", item, cashierToken)
	if created.StatusCode != http.StatusForbidden {
		t.Fatalf("cashier item create status = %d, want %d", created.StatusCode, http.StatusForbidden)
	}
	created.Body.Close()

	created = request(t, app, http.MethodPost, "/api/v1/items", item, managerToken)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("manager item create status = %d, want %d", created.StatusCode, http.StatusCreated)
	}
	created.Body.Close()

	deleted := request(t, app, http.MethodDelete, "/api/v1/items/1", "", managerToken)
	if deleted.StatusCode != http.StatusForbidden {
		t.Fatalf("manager item delete status = %d, want %d", deleted.StatusCode, http.StatusForbidden)
	}
	deleted.Body.Close()

	deleted = request(t, app, http.MethodDelete, "/api/v1/items/1", "", ownerToken)
	if deleted.StatusCode != http.StatusOK {
		t.Fatalf("owner item delete status = %d, want %d", deleted.StatusCode, http.StatusOK)
	}
	deleted.Body.Close()
}

func TestJWTClaimsAndExpiration(t *testing.T) {
	secret := []byte("test-secret-for-jwt-signing")
	token, err := utils.GenerateToken(secret, 42, "owner", "Owner")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := utils.ParseToken(secret, token)
	if err != nil {
		t.Fatalf("parse generated token: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "owner" || claims.Role != "Owner" {
		t.Fatalf("unexpected JWT claims: %#v", claims)
	}
	if time.Until(claims.ExpiresAt.Time) > 24*time.Hour || time.Until(claims.ExpiresAt.Time) <= 23*time.Hour {
		t.Fatalf("unexpected token expiration: %v", claims.ExpiresAt.Time)
	}
	if _, err := utils.ParseToken([]byte("wrong-secret"), token); err == nil {
		t.Fatal("expected token with wrong signing key to fail")
	}
}
