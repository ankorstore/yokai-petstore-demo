-- name: ListOwnerPets :many
SELECT * FROM pets
WHERE owner_id = ?
ORDER BY id;

-- name: CreateOwnerPet :execresult
INSERT INTO pets (
    name, type, owner_id
) VALUES (
   ?, ?, ?
);

-- name: GetOwnerPet :one
SELECT * FROM pets
WHERE owner_id = ?
AND id = ?
LIMIT 1;

-- name: DeleteOwnerPet :exec
DELETE FROM pets
WHERE owner_id = ?
AND id = ?