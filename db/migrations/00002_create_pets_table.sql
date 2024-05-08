-- +goose Up
CREATE TABLE pets (
    id        INTEGER        NOT NULL PRIMARY KEY /*!40101 AUTO_INCREMENT */,
    name      VARCHAR(255)   NOT NULL,
    type      VARCHAR(255)   NOT NULL,
    owner_id  INTEGER,
    FOREIGN KEY (owner_id) REFERENCES owners(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS pets;