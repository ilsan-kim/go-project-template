
CREATE table tasks (
    id          INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    start_date  TIMESTAMP,
    due_date    TIMESTAMP,
    done        BOOLEAN NOT NULL DEFAULT FALSE,
    deleted     BOOLEAN NOT NULL DEFAULT FALSE
);
