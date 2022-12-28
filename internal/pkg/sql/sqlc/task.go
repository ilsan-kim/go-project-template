package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sampleProject/cmd/client"
	"sampleProject/internal/domain"
	"sampleProject/internal/params"
	"sampleProject/internal/pkg/sql/sqlc/mysqlc"
	"sampleProject/internal/service/task"
	"time"
)

type Task struct {
	q  *mysqlc.Queries
	db *sql.DB
	tx *sql.Tx
}

func NewTask() *Task {
	return &Task{
		q:  mysqlc.New(client.SQLClient),
		db: client.SQLClient,
	}
}

func (t *Task) Create(ctx context.Context, input params.CreateParams) (*domain.Task, error) {
	dbConn := t.q
	if _, ok := ctx.Value(task.Transactional{}).(bool); ok {
		dbConn = t.q.WithTx(t.tx)
	}
	result, err := dbConn.InsertTask(ctx, mysqlc.InsertTaskParams{
		Description: input.Description,
		StartDate:   newNullTime(input.StartDate),
		DueDate:     newNullTime(input.DueDate),
	})
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	return &domain.Task{
		ID:          int32(id),
		Description: input.Description,
		StartDate:   input.StartDate,
		DueDate:     input.DueDate,
	}, nil
}

func (t *Task) Get(ctx context.Context, id int) (*domain.Task, error) {
	task, err := t.q.SelectTask(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	// 오래걸리는 척
	time.Sleep(3 * time.Second)

	return &domain.Task{
		ID:          task.ID,
		Description: task.Description,
		Done:        task.Done,
		StartDate:   task.StartDate.Time,
		DueDate:     task.DueDate.Time,
	}, nil
}

func (t *Task) GetList(ctx context.Context) ([]*domain.Task, error) {
	var res []*domain.Task

	tasks, err := t.q.SelectTasks(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		res = append(res, &domain.Task{
			ID:          t.ID,
			Description: t.Description,
			Done:        t.Done,
			StartDate:   t.StartDate.Time,
			DueDate:     t.DueDate.Time,
		})
	}

	return res, nil
}

func (t *Task) Update(ctx context.Context, dao *domain.Task) (*domain.Task, error) {
	dbConn := t.q
	if _, ok := ctx.Value(task.Transactional{}).(bool); ok {
		dbConn = t.q.WithTx(t.tx)
	}

	params := mysqlc.UpdateTaskParams{
		Deleted:     dao.Deleted,
		Description: dao.Description,
		StartDate:   newNullTime(dao.StartDate),
		DueDate:     newNullTime(dao.DueDate),
		Done:        dao.Done,
		ID:          dao.ID,
	}

	_, err := dbConn.UpdateTask(ctx, params)
	if err != nil {
		return nil, err
	}

	return dao, nil
}

func (t *Task) CreateLogTable(ctx context.Context, tableName string) error {
	transactional := false
	if v, ok := ctx.Value(task.Transactional{}).(bool); ok {
		transactional = v
	}

	create := fmt.Sprintf(
		`create table %s (
    		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    		action VARCHAR(255) NOT NULL,
    		created_at TIMESTAMP NOT NULL
		)
	`, tableName)

	switch transactional {
	case false:
		_, err := t.db.ExecContext(ctx, create)
		if err != nil {
			return err
		}
		return nil
	case true:
		_, err := t.tx.ExecContext(ctx, create)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unknown")
	}
}

func (t *Task) Delete(ctx context.Context, id int) error {
	_, err := t.q.DeleteTask(ctx, int32(id))
	if err != nil {
		return err
	}

	// 오래걸리는 척
	fmt.Println("delete process start...")
	time.Sleep(30 * time.Second)
	fmt.Println("delete process done...")
	return nil
}

func (t *Task) Begin(ctx context.Context) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	t.tx = tx
	return nil
}

func (t *Task) Rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func newNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func newNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
