package common

import (
	"slices"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	config    *Config
	scheduler gocron.Scheduler
}

func NewScheduler(config *Config) *Scheduler {
	schedulerInstance := &Scheduler{
		config: config,
	}

	scheduler, newSchedulerErr := gocron.NewScheduler()
	if newSchedulerErr != nil {
		panic(newSchedulerErr)
	}

	if _, newJobErr := scheduler.NewJob(gocron.DurationJob(config.MonitoringInterval), gocron.NewTask(func() {
		mayflyTotalJobs.Set(float64(len(scheduler.Jobs())) - 1)
	}), gocron.WithTags("monitoring")); newJobErr != nil {
		panic(newJobErr)
	}

	schedulerInstance.scheduler = scheduler
	scheduler.Start()

	return schedulerInstance
}

func (s *Scheduler) CreateOrUpdateTask(tag string, date time.Time, task func() error) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		_, updateErr := s.scheduler.Update(job.ID(), gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(date)), gocron.NewTask(task), gocron.WithTags(tag))

		return updateErr
	}

	_, jobErr := s.scheduler.NewJob(gocron.OneTimeJob(
		gocron.OneTimeJobStartDateTime(date)), gocron.NewTask(task), gocron.WithTags(tag))

	return jobErr
}

func (s *Scheduler) DeleteTask(tag string) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		return s.scheduler.RemoveJob(job.ID())
	}

	return nil
}
