package gorm

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sampleProject/cmd/client"
	"sampleProject/internal/domain"
	"sampleProject/internal/params"
	"sampleProject/internal/service/task"
	"time"
)

type Task struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewTask() *Task {
	return &Task{db: client.GormClient}
}

func (t *Task) Create(ctx context.Context, input params.CreateParams) (*domain.Task, error) {
	dbConn := t.db
	if _, ok := ctx.Value(task.Transactional{}).(bool); ok {
		dbConn = t.tx
	}

	task := TaskDAO{
		Description: input.Description,
		StartDate:   input.StartDate,
		DueDate:     input.DueDate,
		Done:        false,
		Deleted:     false,
	}
	res := dbConn.Create(&task)
	if res.RowsAffected == 0 {
		return nil, errors.New("failed to create task data")
	}

	ret := &domain.Task{
		ID:          int32(task.ID),
		Description: task.Description,
		Done:        task.Done,
		StartDate:   task.StartDate,
		DueDate:     task.DueDate,
		Deleted:     task.Deleted,
	}
	return ret, nil
}

func (t *Task) Get(ctx context.Context, id int) (*domain.Task, error) {
	var task TaskDAO

	// 아래처럼 type-safe 하지 않게 쓰일 수도 있다.
	// https://gorm.io/docs/advanced_query.html <- 복잡한 쿼리가 다 type-safe 하지 않음
	//t.db.First(&task, "start_date = ?", "asd")

	res := t.db.First(&task, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}

	ret := &domain.Task{
		ID:          int32(id),
		Description: task.Description,
		Done:        task.Done,
		StartDate:   task.StartDate,
		DueDate:     task.DueDate,
		Deleted:     task.Deleted,
	}
	return ret, nil
}

func (t *Task) GetList(ctx context.Context) ([]*domain.Task, error) {
	var tasks []TaskDAO
	var ret []*domain.Task

	res := t.db.Find(&tasks)
	if res.Error != nil {
		return ret, nil
	}

	for _, t := range tasks {
		ret = append(ret, &domain.Task{
			ID:          int32(t.ID),
			Description: t.Description,
			Done:        t.Done,
			StartDate:   t.StartDate,
			DueDate:     t.DueDate,
			Deleted:     t.Deleted,
		})
	}

	return ret, nil
}

func (t *Task) Update(ctx context.Context, dao *domain.Task) (*domain.Task, error) {
	dbConn := t.db
	if ctx.Value(task.Transactional{}).(bool) {
		dbConn = t.tx
	}

	target := TaskDAO{ID: uint(dao.ID)}

	target.DueDate = dao.DueDate
	target.Done = dao.Done
	target.Deleted = dao.Deleted
	target.StartDate = dao.StartDate
	target.Description = dao.Description
	res := dbConn.Save(target)
	if res.Error != nil {
		return nil, res.Error
	}

	return dao, nil
}

func (t *Task) CreateLogTable(ctx context.Context, tableName string) error {
	dbConn := t.db
	if ctx.Value(task.Transactional{}).(bool) {
		dbConn = t.tx
	}

	query := fmt.Sprintf(`
		create table %s (
    		id int auto_increment not null primary key,
		    created_at timestamp,
		    data varchar(255)
		)
	`, tableName)
	res := dbConn.Exec(query)
	return res.Error
}

func (t *Task) Delete(ctx context.Context, id int) error {
	dbConn := t.db
	if ctx.Value(task.Transactional{}).(bool) {
		dbConn = t.tx
	}

	target := TaskDAO{ID: uint(id)}
	res := dbConn.Delete(target)

	// 오래걸리는 척
	fmt.Println("delete process start...")
	time.Sleep(30 * time.Second)
	fmt.Println("delete process done...")

	return res.Error
}

func (t *Task) Begin(ctx context.Context) error {
	res := t.db.Begin()
	if res.Error != nil {
		return res.Error
	}
	t.tx = res
	return nil
}

func (t *Task) Rollback() error {
	res := t.tx.Rollback()
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (t *Task) Commit() error {
	res := t.tx.Commit()
	if res.Error != nil {
		return res.Error
	}
	return nil
}
