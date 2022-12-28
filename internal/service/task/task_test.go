package task

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sampleProject/internal/domain"
	"sampleProject/internal/params"
	"sampleProject/internal/service/task/taskmock"
	"testing"
	"time"
)

func TestTaskService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_task.NewMockTaskRepository(ctrl)
	cache := mock_task.NewMockTaskCache(ctrl)
	queue := mock_task.NewMockTaskQueue(ctrl)

	svc := NewTaskService(repo, cache, queue)
	t.Run("TODO 리스트 생성", func(t *testing.T) {
		// 여기서 테스트에 사용할 mock 데이터를 mocking된 구조체에 넣어준다.
		input := params.CreateParams{
			Description: "hello",
			StartDate:   time.Now(),
			DueDate:     time.Now().AddDate(0, 0, 1),
		}

		repo.EXPECT().Begin(gomock.Eq(context.WithValue(context.Background(), Transactional{}, true))).Return(nil)
		repo.EXPECT().Rollback().Return(nil)
		repo.EXPECT().Commit().Return(nil)
		repo.EXPECT().Create(gomock.Eq(context.WithValue(context.Background(), Transactional{}, true)), gomock.Eq(input)).
			DoAndReturn(func(_ context.Context, _ params.CreateParams) (*domain.Task, error) {
				return &domain.Task{
					ID:          1,
					Description: input.Description,
					Done:        false,
					StartDate:   input.StartDate,
					DueDate:     input.DueDate,
					Deleted:     false,
				}, nil
			})
		repo.EXPECT().CreateLogTable(context.WithValue(context.Background(), Transactional{}, true), "task_1_log").Return(nil)

		res, err := svc.Create(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, res.Description, input.Description)
	})
}
