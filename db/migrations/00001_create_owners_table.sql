-- +goose Up
CREATE TABLE owners (
    id   INTEGER NOT NULL PRIMARY KEY /*!40101 AUTO_INCREMENT */,
    name VARCHAR(255) NOT NULL,
    bio  VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS owners;