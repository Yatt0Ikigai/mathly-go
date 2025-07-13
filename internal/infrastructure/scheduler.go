//go:generate mockgen -source=scheduler.go -package infrastructure -destination=scheduler_mock.go
//go:generate mockgen -destination=job_mock.go -package=infrastructure github.com/go-co-op/gocron/v2 Job

package infrastructure

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler interface {
	Start()
	Shutdown() error
	NewJob(gocron.JobDefinition, gocron.Task, ...gocron.JobOption) (gocron.Job, error)
	RemoveJob(uuid.UUID) error
}
