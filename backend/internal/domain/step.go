package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)


type StepType string

const (
	StepTypeEmail       StepType = "email"
	StepTypeWait        StepType = "wait"
	StepTypeCondition   StepType = "condition"
	StepTypeAction      StepType = "action"
	StepTypeWebhook     StepType = "webhook"
)

func (st StepType) IsValid() bool {
	switch st {
	case StepTypeEmail, StepTypeWait, StepTypeCondition, StepTypeAction, StepTypeWebhook:
		return true
	default:
		return false
	}
}

type StepOrder struct {
	value int
}

func NewStepOrder(order int) (*StepOrder, error) {
	if order < 1 {
		return nil, ErrInvalidOrder
	}
	return &StepOrder{value: order}, nil
}

func (so StepOrder) Value() int {
	return so.value
}

func (so StepOrder) Next() StepOrder {
	return StepOrder{value: so.value + 1}
}

func (so StepOrder) IsAfter(other StepOrder) bool {
	return so.value > other.value
}

func (so StepOrder) IsBefore(other StepOrder) bool {
	return so.value < other.value
}

type StepConfig struct {
	data map[string]interface{}
}

func NewStepConfig(data map[string]interface{}) *StepConfig {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &StepConfig{data: data}
}

func (sc *StepConfig) Get(key string) (interface{}, bool) {
	value, exists := sc.data[key]
	return value, exists
}

func (sc *StepConfig) Set(key string, value interface{}) {
	sc.data[key] = value
}

func (sc *StepConfig) GetString(key string) (string, bool) {
	if value, exists := sc.data[key]; exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

func (sc *StepConfig) GetInt(key string) (int, bool) {
	if value, exists := sc.data[key]; exists {
		if i, ok := value.(int); ok {
			return i, true
		}
	}
	return 0, false
}

type Step struct {
	ID          uuid.UUID
	Name        string
	Type        StepType
	Order       StepOrder
	Config      *StepConfig
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewStep(name string, stepType StepType, order StepOrder, config *StepConfig, description string, now func() time.Time) (*Step, error) {
	if name == "" {
		return nil, ErrEmptyStepName
	}
	
	if !stepType.IsValid() {
		return nil, ErrInvalidStepType
	}
	
	if config == nil {
		config = NewStepConfig(nil)
	}
	
	if err := validateStepConfig(stepType, config); err != nil {
		return nil, err
	}
	
	nowTime := now()
	return &Step{
		ID:          uuid.New(),
		Name:        name,
		Type:        stepType,
		Order:       order,
		Config:      config,
		Description: description,
		CreatedAt:   nowTime,
		UpdatedAt:   nowTime,
	}, nil
}

func (s *Step) UpdateName(name string, now func() time.Time) error {
	if name == "" {
		return ErrEmptyStepName
	}
	s.Name = name
	s.UpdatedAt = now()
	return nil
}

func (s *Step) UpdateDescription(description string, now func() time.Time) {
	s.Description = description
	s.UpdatedAt = now()
}

func (s *Step) UpdateConfig(config *StepConfig, now func() time.Time) error {
	if config == nil {
		config = NewStepConfig(nil)
	}
	
	if err := validateStepConfig(s.Type, config); err != nil {
		return err
	}
	
	s.Config = config
	s.UpdatedAt = now()
	return nil
}

func (s *Step) IsEmailStep() bool {
	return s.Type == StepTypeEmail
}

func (s *Step) IsWaitStep() bool {
	return s.Type == StepTypeWait
}

func (s *Step) IsConditionStep() bool {
	return s.Type == StepTypeCondition
}

func validateStepConfig(stepType StepType, config *StepConfig) error {
	switch stepType {
	case StepTypeEmail:
		if _, exists := config.GetString("template_id"); !exists {
			return ErrInvalidConfig
		}
	case StepTypeWait:
		if duration, exists := config.GetInt("duration_hours"); !exists || duration < 0 {
			return ErrInvalidConfig
		}
	case StepTypeWebhook:
		if _, exists := config.GetString("url"); !exists {
			return ErrInvalidConfig
		}
	}
	return nil
}