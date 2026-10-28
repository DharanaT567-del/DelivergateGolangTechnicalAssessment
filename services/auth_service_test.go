package services

import (
	"errors"
	"path/filepath"
	"testing"

	"icecream-shop/config"
	"icecream-shop/repository"
)

func TestLoginReturnsDatabaseErrors(t *testing.T) {
	db, err := config.OpenDatabase(filepath.Join(t.TempDir(), "users.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get test database handle: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close test database: %v", err)
	}

	service := NewAuthService(repository.NewUserRepository(db), "test-secret")
	_, err = service.Login("owner", "owner123")
	if err == nil {
		t.Fatal("expected a database error")
	}
	if errors.Is(err, ErrUnauthorized) {
		t.Fatal("database error was reported as invalid credentials")
	}
}
