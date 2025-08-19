package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewWorkflow(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }

	tests := []struct {
		name        string
		workflowName string
		description string
		wantErr     error
	}{
		{
			name:        "valid workflow",
			workflowName: "Test Workflow",
			description: "Test Description",
			wantErr:     nil,
		},
		{
			name:        "empty name",
			workflowName: "",
			description: "Test Description",
			wantErr:     ErrWorkflowNameEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflow, err := NewWorkflow(tt.workflowName, tt.description, now)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewWorkflow() unexpected error = %v", err)
				return
			}

			if workflow.Name != tt.workflowName {
				t.Errorf("workflow.Name = %v, want %v", workflow.Name, tt.workflowName)
			}

			if workflow.Description != tt.description {
				t.Errorf("workflow.Description = %v, want %v", workflow.Description, tt.description)
			}

			if workflow.Status != WorkflowStatusDraft {
				t.Errorf("workflow.Status = %v, want %v", workflow.Status, WorkflowStatusDraft)
			}

			if len(workflow.Steps) != 0 {
				t.Errorf("workflow.Steps length = %v, want 0", len(workflow.Steps))
			}

			if workflow.CreatedAt != now() {
				t.Errorf("workflow.CreatedAt = %v, want %v", workflow.CreatedAt, now())
			}
		})
	}
}

func TestWorkflow_AddStep(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	
	workflow, _ := NewWorkflow("Test Workflow", "Description", now)
	
	order1, _ := NewStepOrder(1)
	step1, _ := NewStep("Step 1", StepTypeEmail, *order1, nil, "Test step", now)
	
	order2, _ := NewStepOrder(2)
	step2, _ := NewStep("Step 2", StepTypeWait, *order2, nil, "Test step 2", now)
	
	orderDuplicate, _ := NewStepOrder(1)
	stepDuplicate, _ := NewStep("Duplicate", StepTypeEmail, *orderDuplicate, nil, "Duplicate step", now)

	tests := []struct {
		name     string
		workflow *Workflow
		step     *Step
		wantErr  error
	}{
		{
			name:     "add first step",
			workflow: workflow,
			step:     step1,
			wantErr:  nil,
		},
		{
			name:     "add second step",
			workflow: workflow,
			step:     step2,
			wantErr:  nil,
		},
		{
			name:     "duplicate order",
			workflow: workflow,
			step:     stepDuplicate,
			wantErr:  ErrDuplicateStepOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.workflow.AddStep(tt.step, now)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("AddStep() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("AddStep() unexpected error = %v", err)
			}
		})
	}
	
	if len(workflow.Steps) != 2 {
		t.Errorf("workflow.Steps length = %v, want 2", len(workflow.Steps))
	}
	
	if workflow.Steps[0].Order.Value() != 1 {
		t.Errorf("first step order = %v, want 1", workflow.Steps[0].Order.Value())
	}
	
	if workflow.Steps[1].Order.Value() != 2 {
		t.Errorf("second step order = %v, want 2", workflow.Steps[1].Order.Value())
	}
}

func TestWorkflow_CannotModifyActiveWorkflow(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	
	workflow, _ := NewWorkflow("Test Workflow", "Description", now)
	
	order1, _ := NewStepOrder(1)
	step1, _ := NewStep("Step 1", StepTypeEmail, nil, nil, "Test step", now)
	step1.Order = *order1
	workflow.AddStep(step1, now)
	
	workflow.Activate(now)

	order2, _ := NewStepOrder(2)
	step2, _ := NewStep("Step 2", StepTypeWait, *order2, nil, "Test step 2", now)

	tests := []struct {
		name     string
		action   func() error
		wantErr  error
	}{
		{
			name: "cannot update name when active",
			action: func() error {
				return workflow.UpdateName("New Name", now)
			},
			wantErr: ErrCannotModifyActiveWorkflow,
		},
		{
			name: "cannot update description when active",
			action: func() error {
				return workflow.UpdateDescription("New Description", now)
			},
			wantErr: ErrCannotModifyActiveWorkflow,
		},
		{
			name: "cannot add step when active",
			action: func() error {
				return workflow.AddStep(step2, now)
			},
			wantErr: ErrCannotModifyActiveWorkflow,
		},
		{
			name: "cannot remove step when active",
			action: func() error {
				return workflow.RemoveStep(step1.ID, now)
			},
			wantErr: ErrCannotModifyActiveWorkflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action()

			if err != tt.wantErr {
				t.Errorf("action() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkflow_StatusTransitions(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	
	workflow, _ := NewWorkflow("Test Workflow", "Description", now)

	tests := []struct {
		name       string
		action     func(*Workflow) error
		wantStatus WorkflowStatus
		wantErr    error
	}{
		{
			name: "cannot activate workflow without steps",
			action: func(w *Workflow) error {
				return w.Activate(now)
			},
			wantStatus: WorkflowStatusDraft,
			wantErr:    ErrStepsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action(workflow)

			if err != tt.wantErr {
				t.Errorf("action() error = %v, wantErr %v", err, tt.wantErr)
			}

			if workflow.Status != tt.wantStatus {
				t.Errorf("workflow.Status = %v, want %v", workflow.Status, tt.wantStatus)
			}
		})
	}
	
	order1, _ := NewStepOrder(1)
	step1, _ := NewStep("Step 1", StepTypeEmail, *order1, nil, "Test step", now)
	workflow.AddStep(step1, now)

	if err := workflow.Activate(now); err != nil {
		t.Errorf("Activate() unexpected error = %v", err)
	}

	if workflow.Status != WorkflowStatusActive {
		t.Errorf("workflow.Status = %v, want %v", workflow.Status, WorkflowStatusActive)
	}

	if err := workflow.Deactivate(now); err != nil {
		t.Errorf("Deactivate() unexpected error = %v", err)
	}

	if workflow.Status != WorkflowStatusInactive {
		t.Errorf("workflow.Status = %v, want %v", workflow.Status, WorkflowStatusInactive)
	}

	if err := workflow.Reactivate(now); err != nil {
		t.Errorf("Reactivate() unexpected error = %v", err)
	}

	if workflow.Status != WorkflowStatusActive {
		t.Errorf("workflow.Status = %v, want %v", workflow.Status, WorkflowStatusActive)
	}
}

func TestWorkflow_GetMethods(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	
	workflow, _ := NewWorkflow("Test Workflow", "Description", now)
	
	order1, _ := NewStepOrder(1)
	step1, _ := NewStep("Step 1", StepTypeEmail, *order1, nil, "Test step", now)
	workflow.AddStep(step1, now)
	
	order2, _ := NewStepOrder(2)
	step2, _ := NewStep("Step 2", StepTypeWait, *order2, nil, "Test step 2", now)
	workflow.AddStep(step2, now)

	t.Run("GetStep", func(t *testing.T) {
		foundStep, err := workflow.GetStep(step1.ID)
		if err != nil {
			t.Errorf("GetStep() unexpected error = %v", err)
		}
		if foundStep.ID != step1.ID {
			t.Errorf("GetStep() returned wrong step")
		}

		_, err = workflow.GetStep(uuid.New())
		if err != ErrStepNotFound {
			t.Errorf("GetStep() error = %v, want %v", err, ErrStepNotFound)
		}
	})

	t.Run("GetStepByOrder", func(t *testing.T) {
		foundStep, err := workflow.GetStepByOrder(*order1)
		if err != nil {
			t.Errorf("GetStepByOrder() unexpected error = %v", err)
		}
		if foundStep.ID != step1.ID {
			t.Errorf("GetStepByOrder() returned wrong step")
		}

		order99, _ := NewStepOrder(99)
		_, err = workflow.GetStepByOrder(*order99)
		if err != ErrStepNotFound {
			t.Errorf("GetStepByOrder() error = %v, want %v", err, ErrStepNotFound)
		}
	})

	t.Run("GetFirstStep", func(t *testing.T) {
		firstStep, err := workflow.GetFirstStep()
		if err != nil {
			t.Errorf("GetFirstStep() unexpected error = %v", err)
		}
		if firstStep.Order.Value() != 1 {
			t.Errorf("GetFirstStep() order = %v, want 1", firstStep.Order.Value())
		}
	})

	t.Run("GetNextStep", func(t *testing.T) {
		nextStep, err := workflow.GetNextStep(*order1)
		if err != nil {
			t.Errorf("GetNextStep() unexpected error = %v", err)
		}
		if nextStep.Order.Value() != 2 {
			t.Errorf("GetNextStep() order = %v, want 2", nextStep.Order.Value())
		}

		_, err = workflow.GetNextStep(*order2)
		if err != ErrStepNotFound {
			t.Errorf("GetNextStep() error = %v, want %v", err, ErrStepNotFound)
		}
	})
}

func TestWorkflow_ValidationAndUtilities(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	
	workflow, _ := NewWorkflow("Test Workflow", "Description", now)

	if workflow.HasSteps() {
		t.Error("HasSteps() should return false for empty workflow")
	}

	if workflow.StepCount() != 0 {
		t.Errorf("StepCount() = %v, want 0", workflow.StepCount())
	}

	if workflow.CanExecute() {
		t.Error("CanExecute() should return false for empty draft workflow")
	}

	order1, _ := NewStepOrder(1)
	step1, _ := NewStep("Step 1", StepTypeEmail, *order1, nil, "Test step", now)
	workflow.AddStep(step1, now)

	if !workflow.HasSteps() {
		t.Error("HasSteps() should return true after adding step")
	}

	if workflow.StepCount() != 1 {
		t.Errorf("StepCount() = %v, want 1", workflow.StepCount())
	}

	if workflow.CanExecute() {
		t.Error("CanExecute() should return false for draft workflow")
	}

	workflow.Activate(now)

	if !workflow.CanExecute() {
		t.Error("CanExecute() should return true for active workflow with steps")
	}

	if workflow.ValidateStepSequence() != nil {
		t.Error("ValidateStepSequence() should return nil for valid sequence")
	}
}