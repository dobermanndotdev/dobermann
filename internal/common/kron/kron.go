package kron

import (
	"context"
	"sync"

	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type Service struct {
	jobs []*Job
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddJob(job *Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(s.jobs))

	for _, job := range s.jobs {
		go func(wg *sync.WaitGroup, job *Job) {
			defer wg.Done()
			err := job.start(ctx)
			if err != nil {
				logs.Error(err)
				return
			}
		}(wg, job)
	}
	wg.Wait()

	return nil
}
