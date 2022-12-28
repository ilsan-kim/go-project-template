package task

import (
	"bytes"
	"context"
	"fmt"
	"sampleProject/internal/common"
	"sampleProject/internal/domain"
	"sampleProject/internal/params"
	"strconv"
	"time"
)

type Transactional struct{}

type TaskRepository interface {
	Create(ctx context.Context, input params.CreateParams) (*domain.Task, error)
	Get(ctx context.Context, id int) (*domain.Task, error)
	GetList(ctx context.Context) ([]*domain.Task, error)
	Update(ctx context.Context, dao *domain.Task) (*domain.Task, error)
	CreateLogTable(ctx context.Context, tableName string) error

	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}

type TaskCache interface {
	Set(ctx context.Context, data *domain.Task) error
	Get(ctx context.Context, id int) (*domain.Task, error)
}

type TaskQueue interface {
	SendMessage(msg string) error
}

type TaskService struct {
	repo  TaskRepository
	cache TaskCache
	queue TaskQueue
}

func NewTaskService(repo TaskRepository, cache TaskCache, queue TaskQueue) *TaskService {
	return &TaskService{
		repo:  repo,
		cache: cache,
		queue: queue,
	}
}

func (t *TaskService) Create(ctx context.Context, input params.CreateParams) (*domain.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// GracefulShutdownJob 잘되나 확인
	dummyFunc := func() {
		fmt.Println("run dummy func for test gracefully shut down")
		time.Sleep(45 * time.Second)
	}
	go common.RunFuncSafely(dummyFunc)

	ctx = context.WithValue(ctx, Transactional{}, true)
	err := t.repo.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer t.repo.Rollback()

	task, err := t.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf("task_%v_log", task.ID)
	err = t.repo.CreateLogTable(ctx, tableName)
	if err != nil {
		return nil, err
	}

	t.repo.Commit()
	return task, nil
}

func (t *TaskService) Get(ctx context.Context, id int) (*domain.Task, error) {
	if task, err := t.cache.Get(ctx, id); err != nil {
		// 캐시에 없다.
		task, err = t.repo.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		err = t.cache.Set(ctx, task)
		return task, nil
	} else {
		// 캐시에 있다.
		return task, nil
	}
}

func (t *TaskService) GetList(ctx context.Context) ([]*domain.Task, error) {
	tasks, err := t.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) MakeTaskListCsv(ctx context.Context) ([]byte, error) {
	// s3에 저장하고 이걸 다운 받는 식으로 만들면 의미는 없지만.. s3를 적용할 수 있다..
	buff := bytes.Buffer{}
	buff.WriteString("id,description,done,startDate,dueDate\n")
	tasks, err := t.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		id := strconv.Itoa(int(t.ID))
		done := "false"

		buff.WriteString(id)
		buff.WriteString(",")

		buff.WriteString(t.Description)
		buff.WriteString(",")

		if t.Done {
			done = "true"
		}
		buff.WriteString(done)
		buff.WriteString(",")

		buff.WriteString(t.StartDate.String())
		buff.WriteString(",")

		buff.WriteString(t.DueDate.String())
		buff.WriteString("\n")
	}

	data := buff.Bytes()
	return data, nil
}

func (t *TaskService) Delete(ctx context.Context, id int) error {
	origTask, err := t.Get(ctx, id)
	if err != nil {
		return err
	}
	target := &domain.Task{
		ID:          int32(id),
		Description: origTask.Description,
		Done:        origTask.Done,
		StartDate:   origTask.StartDate,
		DueDate:     origTask.DueDate,
		Deleted:     true,
	}
	_, err = t.repo.Update(ctx, target)
	if err != nil {
		return err
	}

	strId := strconv.Itoa(id)
	// 하드 딜리트 작업은 여기서 SQS에 넣는식으로 처리
	deleteTaskMsg := fmt.Sprintf("DELETE|%s", strId)
	err = t.queue.SendMessage(deleteTaskMsg)
	if err != nil {
		return err
	}
	return nil
}
