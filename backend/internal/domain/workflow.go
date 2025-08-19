package domain

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)


type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusInactive WorkflowStatus = "inactive"
)

func (ws WorkflowStatus) IsValid() bool {
	switch ws {
	case WorkflowStatusDraft, WorkflowStatusActive, WorkflowStatusInactive:
		return true
	default:
		return false
	}
}

func (ws WorkflowStatus) CanTransitionTo(target WorkflowStatus) bool {
	switch ws {
	case WorkflowStatusDraft:
		return target == WorkflowStatusActive
	case WorkflowStatusActive:
		return target == WorkflowStatusInactive
	case WorkflowStatusInactive:
		return target == WorkflowStatusActive
	default:
		return false
	}
}

type Workflow struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      WorkflowStatus
	Steps       []*Step
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewWorkflow(name, description string, now func() time.Time) (*Workflow, error) {
	if name == "" {
		return nil, ErrWorkflowNameEmpty
	}
	
	nowTime := now()
	return &Workflow{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Status:      WorkflowStatusDraft,
		Steps:       make([]*Step, 0),
		CreatedAt:   nowTime,
		UpdatedAt:   nowTime,
	}, nil
}

func (w *Workflow) UpdateName(name string, now func() time.Time) error {
	if name == "" {
		return ErrWorkflowNameEmpty
	}
	
	if w.Status == WorkflowStatusActive {
		return ErrCannotModifyActiveWorkflow
	}
	
	w.Name = name
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) UpdateDescription(description string, now func() time.Time) error {
	if w.Status == WorkflowStatusActive {
		return ErrCannotModifyActiveWorkflow
	}
	
	w.Description = description
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) AddStep(step *Step, now func() time.Time) error {
	if w.Status == WorkflowStatusActive {
		return ErrCannotModifyActiveWorkflow
	}
	
	if w.hasStepWithOrder(step.Order) {
		return ErrDuplicateStepOrder
	}
	
	w.Steps = append(w.Steps, step)
	w.sortSteps()
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) RemoveStep(stepID uuid.UUID, now func() time.Time) error {
	if w.Status == WorkflowStatusActive {
		return ErrCannotModifyActiveWorkflow
	}
	
	for i, step := range w.Steps {
		if step.ID == stepID {
			w.Steps = append(w.Steps[:i], w.Steps[i+1:]...)
			w.UpdatedAt = now()
			return nil
		}
	}
	
	return ErrStepNotFound
}

func (w *Workflow) GetStep(stepID uuid.UUID) (*Step, error) {
	for _, step := range w.Steps {
		if step.ID == stepID {
			return step, nil
		}
	}
	return nil, ErrStepNotFound
}

func (w *Workflow) GetStepByOrder(order StepOrder) (*Step, error) {
	for _, step := range w.Steps {
		if step.Order.Value() == order.Value() {
			return step, nil
		}
	}
	return nil, ErrStepNotFound
}

func (w *Workflow) GetFirstStep() (*Step, error) {
	if len(w.Steps) == 0 {
		return nil, ErrStepNotFound
	}
	return w.Steps[0], nil
}

func (w *Workflow) GetNextStep(currentOrder StepOrder) (*Step, error) {
	for _, step := range w.Steps {
		if step.Order.IsAfter(currentOrder) {
			return step, nil
		}
	}
	return nil, ErrStepNotFound
}

func (w *Workflow) HasSteps() bool {
	return len(w.Steps) > 0
}

func (w *Workflow) StepCount() int {
	return len(w.Steps)
}

func (w *Workflow) Activate(now func() time.Time) error {
	if !w.Status.CanTransitionTo(WorkflowStatusActive) {
		return ErrInvalidWorkflowTransition
	}
	
	if !w.HasSteps() {
		return ErrStepsRequired
	}
	
	w.Status = WorkflowStatusActive
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) Deactivate(now func() time.Time) error {
	if !w.Status.CanTransitionTo(WorkflowStatusInactive) {
		return ErrInvalidWorkflowTransition
	}
	
	w.Status = WorkflowStatusInactive
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) Reactivate(now func() time.Time) error {
	if w.Status != WorkflowStatusInactive {
		return ErrInvalidWorkflowTransition
	}
	
	if !w.HasSteps() {
		return ErrStepsRequired
	}
	
	w.Status = WorkflowStatusActive
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) IsActive() bool {
	return w.Status == WorkflowStatusActive
}

func (w *Workflow) IsDraft() bool {
	return w.Status == WorkflowStatusDraft
}

func (w *Workflow) IsInactive() bool {
	return w.Status == WorkflowStatusInactive
}

func (w *Workflow) CanExecute() bool {
	return w.IsActive() && w.HasSteps()
}

func (w *Workflow) ReorderSteps(stepOrders map[uuid.UUID]StepOrder, now func() time.Time) error {
	if w.Status == WorkflowStatusActive {
		return ErrCannotModifyActiveWorkflow
	}
	
	orderSet := make(map[int]bool)
	for stepID, newOrder := range stepOrders {
		if orderSet[newOrder.Value()] {
			return ErrDuplicateStepOrder
		}
		orderSet[newOrder.Value()] = true
		
		step, err := w.GetStep(stepID)
		if err != nil {
			return err
		}
		step.Order = newOrder
	}
	
	w.sortSteps()
	w.UpdatedAt = now()
	return nil
}

func (w *Workflow) ValidateStepSequence() error {
	if len(w.Steps) == 0 {
		return nil
	}
	
	orderSet := make(map[int]bool)
	for _, step := range w.Steps {
		if orderSet[step.Order.Value()] {
			return ErrDuplicateStepOrder
		}
		orderSet[step.Order.Value()] = true
	}
	
	return nil
}

func (w *Workflow) hasStepWithOrder(order StepOrder) bool {
	for _, step := range w.Steps {
		if step.Order.Value() == order.Value() {
			return true
		}
	}
	return false
}

func (w *Workflow) sortSteps() {
	sort.Slice(w.Steps, func(i, j int) bool {
		return w.Steps[i].Order.IsBefore(w.Steps[j].Order)
	})
}