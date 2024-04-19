-- name: GetOwner :one
SELECT * FROM owners
WHERE id = ? LIMIT 1;

-- name: ListOwners :many
SELECT * FROM owners
ORDER BY name;

-- name: CreateOwner :execresult
INSERT INTO owners (
    name, bio
) VALUES (
    ?, ?
);

-- name: DeleteOwner :exec
DELETE FROM owners
WHERE id = ?;