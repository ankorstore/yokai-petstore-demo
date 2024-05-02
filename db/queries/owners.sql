-- name: GetOwner :one
SELECT o.*, COUNT(p.id) AS total_pets FROM owners AS o
LEFT JOIN pets AS p ON o.id = p.owner_id
WHERE o.id = $1
GROUP BY o.id, o.name, o.bio
LIMIT 1;

-- name: ListOwners :many
SELECT o.*, COUNT(p.id) AS total_pets FROM owners AS o
LEFT JOIN pets AS p ON o.id = p.owner_id
GROUP BY o.id, o.name, o.bio
ORDER BY o.id;

-- name: CreateOwner :execresult
INSERT INTO owners (
    name, bio
) VALUES (
    $1, $2
);

-- name: DeleteOwner :exec
DELETE FROM owners
WHERE id = $1;