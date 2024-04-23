CREATE TABLE pets (
    id        BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name      text   NOT NULL,
    type      text   NOT NULL,
    owner_id  BIGINT,
    FOREIGN KEY (owner_id) REFERENCES owners(id) ON DELETE CASCADE
);