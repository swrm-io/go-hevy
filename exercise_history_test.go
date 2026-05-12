package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestExerciseHistoryGet(t *testing.T) {
	const id = "11A123F3"
	_, client := newTestServer(t, "/v1/exercise_history/"+id, "exercise_history_get.json", 200)
	entries, err := client.ExerciseHistory.Get(context.Background(), id, nil)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	e := entries[0]
	assert.Equal(t, "5c079430-4d04-4507-9718-e60310665dee", e.WorkoutID)
	assert.Equal(t, "Full Body 3", e.WorkoutTitle)
	assert.Equal(t, id, e.ExerciseTemplateID)
	assert.Equal(t, hevy.SetTypeWarmup, e.SetType)
	require.NotNil(t, e.WeightKg)
	assert.Equal(t, 45.35929094356398, *e.WeightKg)
	require.NotNil(t, e.Reps)
	assert.Equal(t, 12, *e.Reps)
}
