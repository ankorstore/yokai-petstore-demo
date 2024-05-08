// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: pets.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createOwnerPet = `-- name: CreateOwnerPet :execresult
INSERT INTO pets (
    name, type, owner_id
) VALUES (
   ?, ?, ?
)
`

type CreateOwnerPetParams struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	OwnerID sql.NullInt32 `json:"owner_id"`
}

func (q *Queries) CreateOwnerPet(ctx context.Context, arg CreateOwnerPetParams) (sql.Result, error) {
	return q.exec(ctx, q.createOwnerPetStmt, createOwnerPet, arg.Name, arg.Type, arg.OwnerID)
}

const deleteOwnerPet = `-- name: DeleteOwnerPet :exec
DELETE FROM pets
WHERE owner_id = ?
AND id = ?
`

type DeleteOwnerPetParams struct {
	OwnerID sql.NullInt32 `json:"owner_id"`
	ID      int32         `json:"id"`
}

func (q *Queries) DeleteOwnerPet(ctx context.Context, arg DeleteOwnerPetParams) error {
	_, err := q.exec(ctx, q.deleteOwnerPetStmt, deleteOwnerPet, arg.OwnerID, arg.ID)
	return err
}

const getOwnerPet = `-- name: GetOwnerPet :one
SELECT id, name, type, owner_id FROM pets
WHERE owner_id = ?
AND id = ?
LIMIT 1
`

type GetOwnerPetParams struct {
	OwnerID sql.NullInt32 `json:"owner_id"`
	ID      int32         `json:"id"`
}

func (q *Queries) GetOwnerPet(ctx context.Context, arg GetOwnerPetParams) (Pet, error) {
	row := q.queryRow(ctx, q.getOwnerPetStmt, getOwnerPet, arg.OwnerID, arg.ID)
	var i Pet
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.OwnerID,
	)
	return i, err
}

const listOwnerPets = `-- name: ListOwnerPets :many
SELECT id, name, type, owner_id FROM pets
WHERE owner_id = ?
ORDER BY id
`

func (q *Queries) ListOwnerPets(ctx context.Context, ownerID sql.NullInt32) ([]Pet, error) {
	rows, err := q.query(ctx, q.listOwnerPetsStmt, listOwnerPets, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pet
	for rows.Next() {
		var i Pet
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
