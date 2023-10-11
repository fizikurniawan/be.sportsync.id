package mocks

import (
	"sportsync/entities"
)

type FakeTokenUtil struct{}

func (fu *FakeTokenUtil) CreateAccessToken(user *entities.User, secret string, expiry int) (string, error) {
	// Implementasi palsu untuk CreateAccessToken
	return "fake_access_token", nil
}

func (fu *FakeTokenUtil) CreateRefreshToken(user *entities.User, secret string, expiry int) (string, error) {
	// Implementasi palsu untuk CreateRefreshToken
	return "fake_refresh_token", nil
}
