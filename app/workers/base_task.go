package workers

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// NewTaskFromPayload creates a new task from any payload struct
// This helper reduces boilerplate for task creation
//
// Example:
//
//	type MyPayload struct {
//	    UserID int    `json:"user_id"`
//	    Action string `json:"action"`
//	}
//
//	payload := MyPayload{UserID: 123, Action: "process"}
//	task, err := workers.NewTaskFromPayload("my:task", payload)
func NewTaskFromPayload(taskType string, payload any) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(taskType, data), nil
}

// UnmarshalTaskPayload unmarshals task payload into a struct
// This helper reduces boilerplate for payload extraction
//
// Example:
//
//	var payload MyPayload
//	if err := workers.UnmarshalTaskPayload(task, &payload); err != nil {
//	    return err
//	}
func UnmarshalTaskPayload(task *asynq.Task, dest any) error {
	if err := json.Unmarshal(task.Payload(), dest); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return nil
}
