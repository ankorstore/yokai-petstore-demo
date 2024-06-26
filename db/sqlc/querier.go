// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateOwner(ctx context.Context, arg CreateOwnerParams) (sql.Result, error)
	CreateOwnerPet(ctx context.Context, arg CreateOwnerPetParams) (sql.Result, error)
	DeleteOwner(ctx context.Context, id int32) error
	DeleteOwnerPet(ctx context.Context, arg DeleteOwnerPetParams) error
	GetOwner(ctx context.Context, id int32) (GetOwnerRow, error)
	GetOwnerPet(ctx context.Context, arg GetOwnerPetParams) (Pet, error)
	ListOwnerPets(ctx context.Context, ownerID sql.NullInt32) ([]Pet, error)
	ListOwners(ctx context.Context) ([]ListOwnersRow, error)
}

var _ Querier = (*Queries)(nil)
