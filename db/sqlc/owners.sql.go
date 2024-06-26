// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: owners.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createOwner = `-- name: CreateOwner :execresult
INSERT INTO owners (
    name, bio
) VALUES (
    ?, ?
)
`

type CreateOwnerParams struct {
	Name string         `json:"name"`
	Bio  sql.NullString `json:"bio"`
}

func (q *Queries) CreateOwner(ctx context.Context, arg CreateOwnerParams) (sql.Result, error) {
	return q.exec(ctx, q.createOwnerStmt, createOwner, arg.Name, arg.Bio)
}

const deleteOwner = `-- name: DeleteOwner :exec
DELETE FROM owners
WHERE id = ?
`

func (q *Queries) DeleteOwner(ctx context.Context, id int32) error {
	_, err := q.exec(ctx, q.deleteOwnerStmt, deleteOwner, id)
	return err
}

const getOwner = `-- name: GetOwner :one
SELECT o.id, o.name, o.bio, COUNT(p.id) AS total_pets FROM owners AS o
LEFT JOIN pets AS p ON o.id = p.owner_id
WHERE o.id = ?
GROUP BY o.id, o.name, o.bio
LIMIT 1
`

type GetOwnerRow struct {
	ID        int32          `json:"id"`
	Name      string         `json:"name"`
	Bio       sql.NullString `json:"bio"`
	TotalPets int64          `json:"total_pets"`
}

func (q *Queries) GetOwner(ctx context.Context, id int32) (GetOwnerRow, error) {
	row := q.queryRow(ctx, q.getOwnerStmt, getOwner, id)
	var i GetOwnerRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.TotalPets,
	)
	return i, err
}

const listOwners = `-- name: ListOwners :many
SELECT o.id, o.name, o.bio, COUNT(p.id) AS total_pets FROM owners AS o
LEFT JOIN pets AS p ON o.id = p.owner_id
GROUP BY o.id, o.name, o.bio
ORDER BY o.id
`

type ListOwnersRow struct {
	ID        int32          `json:"id"`
	Name      string         `json:"name"`
	Bio       sql.NullString `json:"bio"`
	TotalPets int64          `json:"total_pets"`
}

func (q *Queries) ListOwners(ctx context.Context) ([]ListOwnersRow, error) {
	rows, err := q.query(ctx, q.listOwnersStmt, listOwners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListOwnersRow
	for rows.Next() {
		var i ListOwnersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Bio,
			&i.TotalPets,
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
