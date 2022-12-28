-- name: SelectTasks :many
SELECT
    id,
    description,
    start_date,
    due_date,
    done
FROM
    tasks
WHERE
    deleted = false;

-- name: SelectTask :one
SELECT
    id,
    description,
    start_date,
    due_date,
    done,
    deleted
FROM
    tasks
WHERE
        id = ? AND deleted = false
    LIMIT 1;

-- name: InsertTask :execresult
INSERT INTO tasks (
    description,
    start_date,
    due_date
)
VALUES (
           ?,
           ?,
           ?
       );

-- name: UpdateTask :execresult
UPDATE tasks SET
                 description = ?,
                 start_date  = ?,
                 due_date    = ?,
                 done        = ?,
                 deleted     = ?
WHERE id = ?;

-- name: DeleteTask :execresult
DELETE FROM
    tasks
WHERE
        id = ?;
