package hevy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestWorkoutPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")

		file := fmt.Sprintf("testdata/responses/workout-%s.json", page)
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		_, err = res.Write(data)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	workouts := []hevy.Workout{}
	pager := client.Workouts()
	for x := range pager {
		workouts = append(workouts, x)
	}

	assert.NotEmpty(t, workouts)

	assert.Len(t, workouts, 6)
}
func TestWorkout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/workouts":
			page := req.URL.Query().Get("page")

			file := fmt.Sprintf("testdata/responses/workout-%s.json", page)
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/workouts/count":
			data, err := os.ReadFile("testdata/responses/workout-count.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/workouts/b459cba5-cd6d-463c-abd6-54f8eafcadcb":
			data, err := os.ReadFile("testdata/responses/single-workout.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/workouts/events":
			page := req.URL.Query().Get("page")

			file := fmt.Sprintf("testdata/responses/workout_event-%s.json", page)
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	t.Run("Test All Workouts", func(t *testing.T) {
		workouts, err := client.AllWorkouts()
		for i := range workouts {
			t.Log(workouts[i].ID)
		}
		assert.NoError(t, err)
		assert.NotEmpty(t, workouts)
		assert.Len(t, workouts, 6)
	})

	t.Run("Test Get Workouts", func(t *testing.T) {
		workouts, next, err := client.GetWorkouts(2, 2)
		assert.Equal(t, 3, next)
		assert.NoError(t, err)
		assert.NotEmpty(t, workouts)
		assert.Len(t, workouts, 2)
	})

	t.Run("Test Workout Count", func(t *testing.T) {
		count, err := client.WorkoutCount()

		assert.NoError(t, err)
		assert.Equal(t, 21, count)
	})

	t.Run("Test Single Workout", func(t *testing.T) {
		workoutID, err := uuid.Parse("b459cba5-cd6d-463c-abd6-54f8eafcadcb")

		assert.NoError(t, err)
		workout, err := client.Workout(workoutID)
		assert.NoError(t, err)
		assert.NotEmpty(t, workout)
		assert.Equal(t, workoutID, workout.ID)
		assert.Equal(t, "Morning Workout 💪", workout.Title)
	})

	t.Run("Test Workout Events", func(t *testing.T) {
		since := time.Now()
		events, err := client.WorkoutEvents(since)
		assert.NoError(t, err)

		assert.Len(t, events, 3)
		updated := 0
		deleted := 0

		for _, evnt := range events {
			if evnt.EventType == hevy.DeletedEvent {
				deleted++
			}
			if evnt.EventType == hevy.UpdatedEvent {
				updated++
			}
		}

		assert.Equal(t, 2, updated)
		assert.Equal(t, 1, deleted)
	})
}
