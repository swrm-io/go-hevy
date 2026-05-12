package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestRoutinesList(t *testing.T) {
	client := newTestServer(t, "/v1/routines", "routines_list.json")
	page, err := client.Routines.List(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 18, page.PageCount)
	require.Len(t, page.Routines, 2)

	routine := page.Routines[0]
	assert.Equal(t, "Week 5 to 8 - Day 3", routine.Title)
	require.NotNil(t, routine.FolderID)
	assert.Equal(t, float64(687389), *routine.FolderID)
	assert.Equal(t, 180, routine.Exercises[0].RestSeconds)

	ex1 := routine.Exercises[1]
	require.NotNil(t, ex1.SupersetID)
	assert.Equal(t, 1, *ex1.SupersetID)
}

func TestRoutinesGet(t *testing.T) {
	const id = "0d299174-8660-4b10-918b-e39722d76a13"
	client := newTestServer(t, "/v1/routines/"+id, "routine_get.json")
	routine, err := client.Routines.Get(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, routine.ID)
	assert.Equal(t, "Week 5 to 8 - Day 3", routine.Title)
	require.NotNil(t, routine.FolderID)
	assert.Equal(t, float64(687389), *routine.FolderID)
	require.Len(t, routine.Exercises, 3)
	assert.Equal(t, "Squat (Barbell)", routine.Exercises[0].Title)
	assert.Equal(t, 180, routine.Exercises[0].RestSeconds)
	assert.Nil(t, routine.Exercises[0].SupersetID)
}

func TestRoutinesListInvalidPageSize(t *testing.T) {
	client := newTestServer(t, "/v1/routines", "routines_list.json")
	_, err := client.Routines.List(context.Background(), 1, 11)
	assert.ErrorIs(t, err, hevy.ErrInvalidPageSize)
}

func TestRoutinesLimitExceeded(t *testing.T) {
	client := newErrorServer(t, 403)
	_, err := client.Routines.Create(context.Background(), hevy.RoutineInput{
		Title:     "Test",
		Notes:     "test",
		Exercises: nil,
	})
	assert.ErrorIs(t, err, hevy.ErrRoutineLimitExceeded)
}
