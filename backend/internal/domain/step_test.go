package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStepType_IsValid(t *testing.T) {
	tests := []struct {
		stepType StepType
		want     bool
	}{
		{StepTypeEmail, true},
		{StepTypeWait, true},
		{StepTypeCondition, true},
		{StepTypeAction, true},
		{StepTypeWebhook, true},
		{StepType("invalid"), false},
		{StepType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.stepType), func(t *testing.T) {
			if got := tt.stepType.IsValid(); got != tt.want {
				t.Errorf("StepType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStepOrder(t *testing.T) {
	tests := []struct {
		name    string
		order   int
		wantErr error
	}{
		{
			name:    "valid order",
			order:   1,
			wantErr: nil,
		},
		{
			name:    "valid high order",
			order:   100,
			wantErr: nil,
		},
		{
			name:    "invalid zero order",
			order:   0,
			wantErr: ErrInvalidOrder,
		},
		{
			name:    "invalid negative order",
			order:   -1,
			wantErr: ErrInvalidOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stepOrder, err := NewStepOrder(tt.order)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewStepOrder() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewStepOrder() unexpected error = %v", err)
				return
			}

			if stepOrder.Value() != tt.order {
				t.Errorf("stepOrder.Value() = %v, want %v", stepOrder.Value(), tt.order)
			}
		})
	}
}

func TestStepOrder_Methods(t *testing.T) {
	order1, _ := NewStepOrder(1)
	order2, _ := NewStepOrder(2)
	order3, _ := NewStepOrder(3)

	t.Run("Next", func(t *testing.T) {
		nextOrder := order1.Next()
		if nextOrder.Value() != 2 {
			t.Errorf("Next() = %v, want 2", nextOrder.Value())
		}
	})

	t.Run("IsAfter", func(t *testing.T) {
		if !order2.IsAfter(*order1) {
			t.Error("order2 should be after order1")
		}
		if order1.IsAfter(*order2) {
			t.Error("order1 should not be after order2")
		}
	})

	t.Run("IsBefore", func(t *testing.T) {
		if !order1.IsBefore(*order2) {
			t.Error("order1 should be before order2")
		}
		if order2.IsBefore(*order1) {
			t.Error("order2 should not be before order1")
		}
	})

	t.Run("Comparison consistency", func(t *testing.T) {
		if order2.IsAfter(*order3) || order2.IsBefore(*order1) {
			t.Error("order comparisons are inconsistent")
		}
	})
}

func TestNewStepConfig(t *testing.T) {
	t.Run("with nil data", func(t *testing.T) {
		config := NewStepConfig(nil)
		if config == nil {
			t.Error("NewStepConfig() should not return nil")
		}

		if _, exists := config.Get("nonexistent"); exists {
			t.Error("Get() should return false for nonexistent key")
		}
	})

	t.Run("with data", func(t *testing.T) {
		data := map[string]interface{}{
			"template_id": "test-template",
			"duration":    24,
		}
		config := NewStepConfig(data)

		value, exists := config.Get("template_id")
		if !exists {
			t.Error("Get() should return true for existing key")
		}
		if value != "test-template" {
			t.Errorf("Get() = %v, want 'test-template'", value)
		}
	})
}

func TestStepConfig_TypedGetters(t *testing.T) {
	config := NewStepConfig(map[string]interface{}{
		"string_value": "test",
		"int_value":    42,
		"wrong_type":   123,
	})

	t.Run("GetString", func(t *testing.T) {
		value, exists := config.GetString("string_value")
		if !exists || value != "test" {
			t.Errorf("GetString() = %v, %v, want 'test', true", value, exists)
		}

		_, exists = config.GetString("nonexistent")
		if exists {
			t.Error("GetString() should return false for nonexistent key")
		}

		_, exists = config.GetString("wrong_type")
		if exists {
			t.Error("GetString() should return false for non-string value")
		}
	})

	t.Run("GetInt", func(t *testing.T) {
		value, exists := config.GetInt("int_value")
		if !exists || value != 42 {
			t.Errorf("GetInt() = %v, %v, want 42, true", value, exists)
		}

		_, exists = config.GetInt("nonexistent")
		if exists {
			t.Error("GetInt() should return false for nonexistent key")
		}

		_, exists = config.GetInt("string_value")
		if exists {
			t.Error("GetInt() should return false for non-int value")
		}
	})

	t.Run("Set", func(t *testing.T) {
		config.Set("new_key", "new_value")
		value, exists := config.GetString("new_key")
		if !exists || value != "new_value" {
			t.Errorf("After Set(), GetString() = %v, %v, want 'new_value', true", value, exists)
		}
	})
}

func TestNewStep(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	order, _ := NewStepOrder(1)

	tests := []struct {
		name        string
		stepName    string
		stepType    StepType
		order       StepOrder
		config      *StepConfig
		description string
		wantErr     error
	}{
		{
			name:        "valid email step",
			stepName:    "Send Welcome Email",
			stepType:    StepTypeEmail,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{"template_id": "welcome"}),
			description: "Send welcome email to new users",
			wantErr:     nil,
		},
		{
			name:        "empty name",
			stepName:    "",
			stepType:    StepTypeEmail,
			order:       *order,
			config:      nil,
			description: "",
			wantErr:     ErrEmptyStepName,
		},
		{
			name:        "invalid type",
			stepName:    "Test Step",
			stepType:    StepType("invalid"),
			order:       *order,
			config:      nil,
			description: "",
			wantErr:     ErrInvalidStepType,
		},
		{
			name:        "email step without template_id",
			stepName:    "Invalid Email Step",
			stepType:    StepTypeEmail,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{}),
			description: "",
			wantErr:     ErrInvalidConfig,
		},
		{
			name:        "wait step with valid duration",
			stepName:    "Wait 24 Hours",
			stepType:    StepTypeWait,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{"duration_hours": 24}),
			description: "Wait for 24 hours",
			wantErr:     nil,
		},
		{
			name:        "wait step with invalid duration",
			stepName:    "Invalid Wait Step",
			stepType:    StepTypeWait,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{"duration_hours": -1}),
			description: "",
			wantErr:     ErrInvalidConfig,
		},
		{
			name:        "webhook step with valid URL",
			stepName:    "Call Webhook",
			stepType:    StepTypeWebhook,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{"url": "https://example.com/webhook"}),
			description: "Call external webhook",
			wantErr:     nil,
		},
		{
			name:        "webhook step without URL",
			stepName:    "Invalid Webhook Step",
			stepType:    StepTypeWebhook,
			order:       *order,
			config:      NewStepConfig(map[string]interface{}{}),
			description: "",
			wantErr:     ErrInvalidConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step, err := NewStep(tt.stepName, tt.stepType, tt.order, tt.config, tt.description, now)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewStep() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewStep() unexpected error = %v", err)
				return
			}

			if step.Name != tt.stepName {
				t.Errorf("step.Name = %v, want %v", step.Name, tt.stepName)
			}

			if step.Type != tt.stepType {
				t.Errorf("step.Type = %v, want %v", step.Type, tt.stepType)
			}

			if step.Order.Value() != tt.order.Value() {
				t.Errorf("step.Order = %v, want %v", step.Order.Value(), tt.order.Value())
			}

			if step.Description != tt.description {
				t.Errorf("step.Description = %v, want %v", step.Description, tt.description)
			}

			if step.CreatedAt != now() {
				t.Errorf("step.CreatedAt = %v, want %v", step.CreatedAt, now())
			}

			if step.UpdatedAt != now() {
				t.Errorf("step.UpdatedAt = %v, want %v", step.UpdatedAt, now())
			}
		})
	}
}

func TestStep_UpdateMethods(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	later := func() time.Time { return time.Date(2025, 8, 19, 13, 0, 0, 0, time.UTC) }
	
	order, _ := NewStepOrder(1)
	step, _ := NewStep("Test Step", StepTypeEmail, *order, 
		NewStepConfig(map[string]interface{}{"template_id": "test"}), "Description", now)

	t.Run("UpdateName", func(t *testing.T) {
		err := step.UpdateName("Updated Name", later)
		if err != nil {
			t.Errorf("UpdateName() unexpected error = %v", err)
		}

		if step.Name != "Updated Name" {
			t.Errorf("step.Name = %v, want 'Updated Name'", step.Name)
		}

		if step.UpdatedAt != later() {
			t.Errorf("step.UpdatedAt was not updated")
		}

		err = step.UpdateName("", later)
		if err != ErrEmptyStepName {
			t.Errorf("UpdateName() error = %v, want %v", err, ErrEmptyStepName)
		}
	})

	t.Run("UpdateDescription", func(t *testing.T) {
		step.UpdateDescription("Updated Description", later)

		if step.Description != "Updated Description" {
			t.Errorf("step.Description = %v, want 'Updated Description'", step.Description)
		}

		if step.UpdatedAt != later() {
			t.Errorf("step.UpdatedAt was not updated")
		}
	})

	t.Run("UpdateConfig", func(t *testing.T) {
		newConfig := NewStepConfig(map[string]interface{}{"template_id": "new-template"})
		err := step.UpdateConfig(newConfig, later)
		if err != nil {
			t.Errorf("UpdateConfig() unexpected error = %v", err)
		}

		templateID, exists := step.Config.GetString("template_id")
		if !exists || templateID != "new-template" {
			t.Errorf("config was not updated correctly")
		}

		if step.UpdatedAt != later() {
			t.Errorf("step.UpdatedAt was not updated")
		}

		invalidConfig := NewStepConfig(map[string]interface{}{})
		err = step.UpdateConfig(invalidConfig, later)
		if err != ErrInvalidConfig {
			t.Errorf("UpdateConfig() error = %v, want %v", err, ErrInvalidConfig)
		}
	})
}

func TestStep_TypeCheckers(t *testing.T) {
	now := func() time.Time { return time.Date(2025, 8, 19, 12, 0, 0, 0, time.UTC) }
	order, _ := NewStepOrder(1)

	tests := []struct {
		stepType StepType
		checker  func(*Step) bool
		want     bool
	}{
		{StepTypeEmail, (*Step).IsEmailStep, true},
		{StepTypeWait, (*Step).IsWaitStep, true},
		{StepTypeCondition, (*Step).IsConditionStep, true},
		{StepTypeEmail, (*Step).IsWaitStep, false},
		{StepTypeWait, (*Step).IsEmailStep, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.stepType), func(t *testing.T) {
			step, _ := NewStep("Test Step", tt.stepType, *order, nil, "", now)

			if got := tt.checker(step); got != tt.want {
				t.Errorf("type checker = %v, want %v", got, tt.want)
			}
		})
	}
}