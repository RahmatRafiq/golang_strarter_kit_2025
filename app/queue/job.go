package queue

import (
	"encoding/json"
	"fmt"
	"time"
)

// JobInterface defines the interface that all jobs must implement
type JobInterface interface {
	Handle() error
	GetName() string
	GetData() map[string]interface{}
	SetData(data map[string]interface{})
}

// BaseJob provides a basic implementation of JobInterface
type BaseJob struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

func (j *BaseJob) Handle() error {
	// Default implementation - should be overridden by specific jobs
	return nil
}

func (j *BaseJob) GetName() string {
	return j.Name
}

func (j *BaseJob) GetData() map[string]interface{} {
	return j.Data
}

func (j *BaseJob) SetData(data map[string]interface{}) {
	j.Data = data
}

// QueueJob represents a job in the queue
type QueueJob struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:"data"`
	Attempts  int                    `json:"attempts"`
	MaxTries  int                    `json:"max_tries"`
	CreatedAt time.Time              `json:"created_at"`
	Delay     time.Duration          `json:"delay"`
}

// QueueInterface defines the interface for queue operations
type QueueInterface interface {
	Push(job JobInterface, delay ...time.Duration) error
	Pop(queueName string) (*QueueJob, error)
	Size(queueName string) (int64, error)
	Clear(queueName string) error
	Failed(job *QueueJob, err error) error
	Retry(job *QueueJob) error
}

// JobRegistry stores registered job types
var JobRegistry = make(map[string]func() JobInterface)

// RegisterJob registers a job type
func RegisterJob(name string, factory func() JobInterface) {
	JobRegistry[name] = factory
}

// CreateJob creates a job instance from registry
func CreateJob(name string, data map[string]interface{}) (JobInterface, error) {
	factory, exists := JobRegistry[name]
	if !exists {
		return nil, fmt.Errorf("job type '%s' not registered", name)
	}

	job := factory()
	job.SetData(data)
	return job, nil
}

// SerializeJob converts a job to JSON
func SerializeJob(job JobInterface) ([]byte, error) {
	queueJob := &QueueJob{
		Name:      job.GetName(),
		Data:      job.GetData(),
		CreatedAt: time.Now(),
		MaxTries:  3,
	}

	return json.Marshal(queueJob)
}

// DeserializeJob converts JSON to a job
func DeserializeJob(data []byte) (*QueueJob, error) {
	var queueJob QueueJob
	err := json.Unmarshal(data, &queueJob)
	return &queueJob, err
}
