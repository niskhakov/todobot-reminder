package jobqueue

import (
	"context"
	"time"
)

type JobID string

type Job struct {
	F  func(ctx context.Context)
	ID JobID
}

type Scheduler interface {
	Add(ctx context.Context, j Job, t time.Duration)
	AddAt(ctx context.Context, j Job, t time.Time)
	Stop(j JobID)
	StopAll()
}
