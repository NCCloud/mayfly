package common

import (
	"slices"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler interface {
	CreateOrUpdateOneTimeTask(tag string, at time.Time, task func() error) error
	CreateOrUpdateRecurringTask(tag string, cron string, task func() error) error
	DeleteTask(tag string) error
	GetTaskNextRun(tag string) string
}

type scheduler struct {
	config    *Config
	scheduler gocron.Scheduler
}

func NewScheduler(config *Config) Scheduler {
	schedulerInstance := &scheduler{
		config: config,
	}

	cronScheduler, newSchedulerErr := gocron.NewScheduler()
	if newSchedulerErr != nil {
		panic(newSchedulerErr)
	}

	if _, newJobErr := cronScheduler.NewJob(gocron.DurationJob(config.MonitoringInterval), gocron.NewTask(func() {
		mayflyTotalJobs.Set(float64(len(cronScheduler.Jobs())) - 1)
	}), gocron.WithTags("monitoring")); newJobErr != nil {
		panic(newJobErr)
	}

	schedulerInstance.scheduler = cronScheduler
	cronScheduler.Start()

	return schedulerInstance
}

func (s *scheduler) CreateOrUpdateOneTimeTask(tag string, date time.Time, task func() error) error {
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

func (s *scheduler) CreateOrUpdateRecurringTask(tag string, cron string, task func() error) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		_, updateErr := s.scheduler.Update(job.ID(), gocron.CronJob(
			cron, true), gocron.NewTask(task), gocron.WithTags(tag))

		return updateErr
	}

	_, jobErr := s.scheduler.NewJob(gocron.CronJob(
		cron, true), gocron.NewTask(task), gocron.WithTags(tag))

	return jobErr
}

func (s *scheduler) GetTaskNextRun(tag string) string {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		nextRun, nextRunErr := job.NextRun()
		if nextRunErr == nil {
			return nextRun.Format(time.RFC3339)
		}
	}

	return ""
}

func (s *scheduler) DeleteTask(tag string) error {
	job := pie.Of(s.scheduler.Jobs()).Filter(func(job gocron.Job) bool {
		return slices.Contains(job.Tags(), tag)
	}).First()

	if job != nil {
		return s.scheduler.RemoveJob(job.ID())
	}

	return nil
}
