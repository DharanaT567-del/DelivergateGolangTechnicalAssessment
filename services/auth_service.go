package services

import (
	"errors"
	"fmt"

	"icecream-shop/repository"
	"icecream-shop/utils"
)

var ErrUnauthorized = errors.New("invalid username or password")

type AuthService struct {
	users  *repository.UserRepository
	secret []byte
}

func NewAuthService(users *repository.UserRepository, secret string) *AuthService {
	return &AuthService{users: users, secret: []byte(secret)}
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.users.FindByUsername(username)
	if err != nil || !utils.CheckPassword(user.Password, password) {
		return "", ErrUnauthorized
	}
	token, err := utils.GenerateToken(s.secret, user.ID, user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return token, nil
}
