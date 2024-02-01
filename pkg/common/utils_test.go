package common

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResolveSchedule_Date(t *testing.T) {
	// given
	date := gofakeit.Date()

	// when
	schedule, scheduleErr := ResolveSchedule(metav1.Time{}, date.String())

	// then
	assert.Nil(t, scheduleErr)
	assert.Equal(t, date.Year(), schedule.Year())
	assert.Equal(t, date.Month(), schedule.Month())
	assert.Equal(t, date.Day(), schedule.Day())
	assert.Equal(t, date.Hour(), schedule.Hour())
	assert.Equal(t, date.Minute(), schedule.Minute())
	assert.Equal(t, date.Second(), schedule.Second())
}

func TestResolveSchedule_Duration(t *testing.T) {
	// given
	date := metav1.Time{Time: gofakeit.Date()}
	duration := 60 * time.Minute

	// when
	schedule, scheduleErr := ResolveSchedule(date, duration.String())

	// then
	assert.Nil(t, scheduleErr)
	assert.Equal(t, date.Year(), schedule.Year())
	assert.Equal(t, date.Month(), schedule.Month())
	assert.Equal(t, date.Day(), schedule.Day())
	assert.Equal(t, date.Hour(), schedule.Hour()-1)
	assert.Equal(t, date.Minute(), schedule.Minute())
	assert.Equal(t, date.Second(), schedule.Second())
}

func TestNewResourceInstance(t *testing.T) {
	// given
	apiVersionKind := "v1;Secret"

	// when
	resourceInstance := NewResourceInstance(apiVersionKind)

	// then
	assert.Equal(t, apiVersionKind, resourceInstance.GetAPIVersion()+";"+resourceInstance.GetKind())
}

func TestNewResourceInstanceList(t *testing.T) {
	// given
	apiVersionKind := "v1;Secret"

	// when
	resourceInstanceList := NewResourceInstanceList(apiVersionKind)

	// then
	assert.Equal(t, apiVersionKind+"List", resourceInstanceList.GetAPIVersion()+";"+resourceInstanceList.GetKind())
}
