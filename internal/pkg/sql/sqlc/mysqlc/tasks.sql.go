// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: tasks.sql

package mysqlc

import (
	"context"
	"database/sql"
)

const deleteTask = `-- name: DeleteTask :execresult
DELETE FROM
    tasks
WHERE
        id = ?
`

func (q *Queries) DeleteTask(ctx context.Context, id int32) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteTask, id)
}

const getCoupon = `-- name: GetCoupon :one
select cno, ifnull(name, 'qwe'), ifnull(owner,'ttq') from coupon where owner = ?
`

type GetCouponRow struct {
	Cno      string      `json:"cno"`
	Ifnull   interface{} `json:"ifnull"`
	Ifnull_2 interface{} `json:"ifnull_2"`
}

func (q *Queries) GetCoupon(ctx context.Context, owner sql.NullString) (GetCouponRow, error) {
	row := q.db.QueryRowContext(ctx, getCoupon, owner)
	var i GetCouponRow
	err := row.Scan(&i.Cno, &i.Ifnull, &i.Ifnull_2)
	return i, err
}

const insertTask = `-- name: InsertTask :execresult
INSERT INTO tasks (
    description,
    start_date,
    due_date
)
VALUES (
           ?,
           ?,
           ?
       )
`

type InsertTaskParams struct {
	Description string       `json:"description"`
	StartDate   sql.NullTime `json:"start_date"`
	DueDate     sql.NullTime `json:"due_date"`
}

func (q *Queries) InsertTask(ctx context.Context, arg InsertTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertTask, arg.Description, arg.StartDate, arg.DueDate)
}

const selectTask = `-- name: SelectTask :one
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
    LIMIT 1
`

func (q *Queries) SelectTask(ctx context.Context, id int32) (Task, error) {
	row := q.db.QueryRowContext(ctx, selectTask, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.StartDate,
		&i.DueDate,
		&i.Done,
		&i.Deleted,
	)
	return i, err
}

const selectTasks = `-- name: SelectTasks :many
SELECT
    id,
    description,
    start_date,
    due_date,
    done
FROM
    tasks
WHERE
    deleted = false
`

type SelectTasksRow struct {
	ID          int32        `json:"id"`
	Description string       `json:"description"`
	StartDate   sql.NullTime `json:"start_date"`
	DueDate     sql.NullTime `json:"due_date"`
	Done        bool         `json:"done"`
}

func (q *Queries) SelectTasks(ctx context.Context) ([]SelectTasksRow, error) {
	rows, err := q.db.QueryContext(ctx, selectTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectTasksRow
	for rows.Next() {
		var i SelectTasksRow
		if err := rows.Scan(
			&i.ID,
			&i.Description,
			&i.StartDate,
			&i.DueDate,
			&i.Done,
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

const updateTask = `-- name: UpdateTask :execresult
UPDATE tasks SET
                 description = ?,
                 start_date  = ?,
                 due_date    = ?,
                 done        = ?,
                 deleted     = ?
WHERE id = ?
`

type UpdateTaskParams struct {
	Description string       `json:"description"`
	StartDate   sql.NullTime `json:"start_date"`
	DueDate     sql.NullTime `json:"due_date"`
	Done        bool         `json:"done"`
	Deleted     bool         `json:"deleted"`
	ID          int32        `json:"id"`
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTask,
		arg.Description,
		arg.StartDate,
		arg.DueDate,
		arg.Done,
		arg.Deleted,
		arg.ID,
	)
}