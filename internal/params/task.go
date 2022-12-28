package params

import (
	"errors"
	"time"
)

type CreateParams struct {
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	DueDate     time.Time `json:"due_date"`
}

func (t *CreateParams) Validate() error {
	if t.Description == "" || t.StartDate.IsZero() || t.DueDate.IsZero() {
		return errors.New("field required")
	}
	if t.StartDate.After(t.DueDate) {
		return errors.New("start date should be before end date")
	}
	return nil
}
