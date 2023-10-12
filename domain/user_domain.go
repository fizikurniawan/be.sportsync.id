package domain

import (
	"context"
	"sportsync/entities"
)

const CollectionUser = "users"

type UserRepository interface {
	Create(c context.Context, user *entities.User) error
	Fetch(c context.Context) ([]entities.User, error)
	GetByEmail(c context.Context, email string) (entities.User, error)
	GetByID(c context.Context, id string) (entities.User, error)
}
