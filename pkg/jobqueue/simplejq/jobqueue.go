package simplejq

import (
	"context"
	"sync"
	"time"

	"github.com/niskhakov/todobot-reminder/pkg/jobqueue"
)

type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations map[jobqueue.JobID]context.CancelFunc
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make(map[jobqueue.JobID]context.CancelFunc),
	}
}

func (s *Scheduler) Add(ctx context.Context, j jobqueue.Job, interval time.Duration) {
	ctx = s.add(ctx, &j)

	go s.process(ctx, j, interval)
}

func (s *Scheduler) AddAt(ctx context.Context, j jobqueue.Job, t time.Time) {
	ctx = s.add(ctx, &j)

	go s.processOnce(ctx, j, t)
}

func (s *Scheduler) StopAll() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) Stop(id jobqueue.JobID) {
	if v, ok := s.cancellations[id]; ok {
		v()
		delete(s.cancellations, id)
	}
}

func (s *Scheduler) add(ctx context.Context, j *jobqueue.Job) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations[j.ID] = cancel

	s.wg.Add(1)
	return ctx
}

func (s *Scheduler) process(ctx context.Context, j jobqueue.Job, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// First Run
	j.F(ctx)
	for {
		select {
		case <-ticker.C:
			j.F(ctx)
		case <-ctx.Done():
			s.wg.Done()
			return
		}
	}
}

func (s *Scheduler) processOnce(ctx context.Context, j jobqueue.Job, t time.Time) {
	defer s.wg.Done()

	dur := time.Until(t)
	if dur < 0 {
		// Time in the past
		return
	}
	timer := time.NewTimer(dur)

	select {
	case <-timer.C:
		j.F(ctx)
	case <-ctx.Done():
		return
	}

}
