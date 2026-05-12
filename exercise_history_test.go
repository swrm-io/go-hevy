package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestExerciseHistoryGet(t *testing.T) {
	const exerciseID = "11A123F3"
	client := newTestServer(t, "/v1/exercise_history/"+exerciseID, "exercise_history_get.json")
	entries, err := client.ExerciseHistory.Get(context.Background(), exerciseID, nil)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	entry := entries[0]
	assert.Equal(t, "5c079430-4d04-4507-9718-e60310665dee", entry.WorkoutID)
	assert.Equal(t, "Full Body 3", entry.WorkoutTitle)
	assert.Equal(t, exerciseID, entry.ExerciseTemplateID)
	assert.Equal(t, hevy.SetTypeWarmup, entry.SetType)
	require.NotNil(t, entry.WeightKg)
	assert.Equal(t, 45.35929094356398, *entry.WeightKg)
	require.NotNil(t, entry.Reps)
	assert.Equal(t, 12, *entry.Reps)
}
