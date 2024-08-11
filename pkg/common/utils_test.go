package common

import (
	"math/rand"
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
	schedule, scheduleErr := ResolveOneTimeSchedule(metav1.Time{}, date.String())

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
	currentDate := metav1.Time{Time: gofakeit.Date()}
	duration := time.Duration(rand.Int31n(100000)) * time.Second
	expected := currentDate.Add(duration)

	// when
	schedule, scheduleErr := ResolveOneTimeSchedule(currentDate, duration.String())

	// then
	assert.Nil(t, scheduleErr)
	assert.Equal(t, expected.Year(), schedule.Year())
	assert.Equal(t, expected.Month(), schedule.Month())
	assert.Equal(t, expected.Day(), schedule.Day())
	assert.Equal(t, expected.Hour(), schedule.Hour())
	assert.Equal(t, expected.Minute(), schedule.Minute())
	assert.Equal(t, expected.Second(), schedule.Second())
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
