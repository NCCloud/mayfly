package common

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_New(t *testing.T) {
	// given
	config := NewConfig()

	// when
	schedulerInstance := NewScheduler(config)

	// then
	assert.NotNil(t, schedulerInstance)
	assert.IsType(t, schedulerInstance, &scheduler{})
}

func TestScheduler_CreateOrUpdateOneTimeTask(t *testing.T) {
	// given
	config := NewConfig()
	schedulerInstance := NewScheduler(config)

	// when
	createOrUpdateTaskErr := schedulerInstance.
		CreateOrUpdateOneTimeTask("monitoring", gofakeit.FutureDate().Add(1*time.Hour), func() error {
			return nil
		})

	// then
	assert.Nil(t, createOrUpdateTaskErr)
}

func TestScheduler_CreateOrUpdateRecurringTask(t *testing.T) {
	// given
	config := NewConfig()
	schedulerInstance := NewScheduler(config)

	// when
	createOrUpdateTaskErr := schedulerInstance.
		CreateOrUpdateRecurringTask("monitoring", "* * * * * *", func() error {
			return nil
		})

	// then
	assert.Nil(t, createOrUpdateTaskErr)
}

func TestScheduler_DeleteTask(t *testing.T) {
	// given
	config := NewConfig()
	schedulerInstance := NewScheduler(config)

	// when
	deleteTaskErr := schedulerInstance.DeleteTask("monitoring")

	// then
	assert.Nil(t, deleteTaskErr)
}
