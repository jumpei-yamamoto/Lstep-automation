package domain

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowRepository interface {
	Save(ctx context.Context, workflow *Workflow) error
	FindByID(ctx context.Context, id uuid.UUID) (*Workflow, error)
	FindByName(ctx context.Context, name string) (*Workflow, error)
	FindActiveWorkflows(ctx context.Context) ([]*Workflow, error)
	FindAll(ctx context.Context) ([]*Workflow, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}