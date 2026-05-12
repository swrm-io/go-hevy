package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestWorkoutsList(t *testing.T) {
	_, client := newTestServer(t, "/v1/workouts", "workouts_list.json", 200)
	page, err := client.Workouts.List(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 109, page.PageCount)
	require.Len(t, page.Workouts, 2)

	w := page.Workouts[0]
	assert.Equal(t, "5c079430-4d04-4507-9718-e60310665dee", w.ID)
	assert.Equal(t, "Full Body 3", w.Title)
	require.Len(t, w.Exercises, 2)
	assert.Nil(t, w.Exercises[0].SupersetID)

	set0 := w.Exercises[0].Sets[0]
	assert.Equal(t, hevy.SetTypeWarmup, set0.Type)
	require.NotNil(t, set0.WeightKg)
	assert.Equal(t, 45.35929094356398, *set0.WeightKg)
}

func TestWorkoutsGet(t *testing.T) {
	const id = "5c079430-4d04-4507-9718-e60310665dee"
	_, client := newTestServer(t, "/v1/workouts/"+id, "workout_get.json", 200)
	w, err := client.Workouts.Get(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, w.ID)
	assert.Equal(t, "Full Body 3", w.Title)
	assert.Equal(t, "c6425b18-2c45-422a-99bc-9d6dd0ae5985", w.RoutineID)
	require.Len(t, w.Exercises, 5)

	set0 := w.Exercises[0].Sets[0]
	assert.Equal(t, hevy.SetTypeWarmup, set0.Type)
	require.NotNil(t, set0.Reps)
	assert.Equal(t, 12, *set0.Reps)
}

func TestWorkoutsCount(t *testing.T) {
	_, client := newTestServer(t, "/v1/workouts/count", "workouts_count.json", 200)
	count, err := client.Workouts.Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 218, count)
}

func TestWorkoutsEvents(t *testing.T) {
	_, client := newTestServer(t, "/v1/workouts/events", "workouts_events.json", 200)
	page, err := client.Workouts.Events(context.Background(), 1, 2, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, page.Page)
	assert.Equal(t, 2, page.PageCount)
	require.Len(t, page.Events, 2)

	e := page.Events[0]
	assert.Equal(t, hevy.WorkoutEventUpdated, e.Type)
	require.NotNil(t, e.Workout)
	assert.Equal(t, "5c079430-4d04-4507-9718-e60310665dee", e.Workout.ID)
	assert.Equal(t, "Full Body 3", e.Workout.Title)
}

func TestWorkoutsListInvalidPageSize(t *testing.T) {
	_, client := newTestServer(t, "/v1/workouts", "workouts_list.json", 200)
	_, err := client.Workouts.List(context.Background(), 1, 99)
	assert.ErrorIs(t, err, hevy.ErrInvalidPageSize)
}

func TestWorkoutsListNoMorePages(t *testing.T) {
	_, client := newErrorServer(t, 404)
	_, err := client.Workouts.List(context.Background(), 999, 10)
	assert.ErrorIs(t, err, hevy.ErrNoMorePages)
}

func TestWorkoutsGetNotFound(t *testing.T) {
	_, client := newErrorServer(t, 404)
	_, err := client.Workouts.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, hevy.ErrNotFound)
}

func TestWorkoutsUnauthorized(t *testing.T) {
	_, client := newErrorServer(t, 401)
	_, err := client.Workouts.Count(context.Background())
	assert.ErrorIs(t, err, hevy.ErrUnauthorized)
}
