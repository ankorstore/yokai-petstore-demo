-- name: ListOwnerPets :many
SELECT * FROM pets
WHERE owner_id = $1
ORDER BY id;

-- name: CreateOwnerPet :execresult
INSERT INTO pets (
    name, type, owner_id
) VALUES (
   $1, $2, $3
);

-- name: GetOwnerPet :one
SELECT * FROM pets
WHERE owner_id = $1
AND id = $2
LIMIT 1;

-- name: DeleteOwnerPet :exec
DELETE FROM pets
WHERE owner_id = $1
AND id = $2;