package hevy_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestUnmarshal(t *testing.T) {
	t.Run("Workout", func(t *testing.T) {
		workout := hevy.Workout{}

		data, err := os.ReadFile("testdata/base/workout.json")
		assert.NoError(t, err, "error reading testdata/base/workout.json")
		err = json.Unmarshal(data, &workout)
		assert.NoError(t, err, "error unmarshalling testdata/base/workout.json")

		assert.NotEmpty(t, workout)
		assert.IsType(t, uuid.UUID{}, workout.ID)

		assert.NotEmpty(t, workout.StartTime)
		assert.NotEmpty(t, workout.EndTime)
		assert.NotEmpty(t, workout.Exercises)

		assert.Equal(t, float64(1564.8955375529572), workout.VolumeKG)
		assert.Equal(t, float64(3450), workout.VolumeLB)

	})
	t.Run("Exercise", func(t *testing.T) {
		exercise := hevy.Exercise{}

		data, err := os.ReadFile("testdata/base/exercise.json")
		assert.NoError(t, err, "error reading testdata/base/exercise.json")
		err = json.Unmarshal(data, &exercise)
		assert.NoError(t, err, "error unmarshalling testdata/base/exercise.json")

		assert.NotEmpty(t, exercise)
		assert.NotEmpty(t, exercise.Sets)

		assert.Equal(t, float64(1564.8955375529572), exercise.VolumeKG)
		assert.Equal(t, float64(3450), exercise.VolumeLB)

	})
	t.Run("Set", func(t *testing.T) {
		set := hevy.Set{}

		data, err := os.ReadFile("testdata/base/set.json")
		assert.NoError(t, err, "error reading testdata/base/set.json")
		err = json.Unmarshal(data, &set)
		assert.NoError(t, err, "error unmarshalling testdata/base/exercise.json")

		assert.NotEmpty(t, set)
		assert.Equal(t, hevy.WarmupSet, set.SetType)

		assert.Equal(t, float64(29.483539113316585), set.WeightKG)
		assert.Equal(t, float64(65), set.WeightLB)

		assert.Equal(t, float64(650), set.VolumeLB)
		assert.Equal(t, float64(294.83539113316584), set.VolumeKG)
	})
	t.Run("Routine", func(t *testing.T) {
		routine := hevy.Routine{}

		data, err := os.ReadFile("testdata/base/routine.json")
		assert.NoError(t, err, "error reading testdata/base/routine.json")
		err = json.Unmarshal(data, &routine)
		assert.NoError(t, err, "error unmarshalling testdata/base/routine.json")

		assert.NotEmpty(t, routine)
		assert.IsType(t, uuid.UUID{}, routine.ID)
		assert.NotEmpty(t, routine.CreatedAt)
		assert.NotEmpty(t, routine.UpdatedAt)
		assert.NotEmpty(t, routine.Exercises)

	})

	t.Run("Workout Event Delete", func(t *testing.T) {
		event := hevy.Event{}

		data, err := os.ReadFile("testdata/base/event_delete.json")
		assert.NoError(t, err, "error reading testdata/base/event_delete.json")
		err = json.Unmarshal(data, &event)
		assert.NoError(t, err, "error unmarshalling testdata/base/event_delete.json")

		assert.NotEmpty(t, event)
		assert.Equal(t, hevy.DeletedEvent, event.EventType)
		assert.IsType(t, uuid.UUID{}, event.ID)
		assert.NotEmpty(t, event.ID)
		assert.Empty(t, event.Workout)
	})

	t.Run("Workup Event Update", func(t *testing.T) {
		event := hevy.Event{}
		data, err := os.ReadFile("testdata/base/event_update.json")
		assert.NoError(t, err, "error reading testdata/base/event_update.json")
		err = json.Unmarshal(data, &event)
		assert.NoError(t, err, "error unmaarshalling testdata/base/event_update.json")

		assert.NotEmpty(t, event)
		assert.Equal(t, hevy.UpdatedEvent, event.EventType)
		assert.Empty(t, event.DeletedAt)
		assert.Empty(t, event.ID)
		assert.NotEmpty(t, event.Workout)
	})

	t.Run("Body Measurement", func(t *testing.T) {
		measurement := hevy.BodyMeasurement{}

		data, err := os.ReadFile("testdata/base/bodymeasurement.json")
		assert.NoError(t, err, "error reading testdata/base/bodymeasurement.json")
		err = json.Unmarshal(data, &measurement)
		assert.NoError(t, err, "error unmarshalling testdata/base/bodymeasurement.json")

		assert.NotEmpty(t, measurement)
		assert.IsType(t, int(0), measurement.ID)
		assert.NotEmpty(t, measurement.CreatedAt)

	})
}
