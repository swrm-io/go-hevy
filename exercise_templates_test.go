package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestExerciseTemplatesList(t *testing.T) {
	_, client := newTestServer(t, "/v1/exercise_templates", "exercise_templates_list.json", 200)
	page, err := client.ExerciseTemplates.List(context.Background(), 1, 3)
	require.NoError(t, err)
	assert.Equal(t, 146, page.PageCount)
	require.Len(t, page.ExerciseTemplates, 3)

	tmpl := page.ExerciseTemplates[0]
	assert.Equal(t, "3BC06AD3", tmpl.ID)
	assert.Equal(t, "21s Bicep Curl", tmpl.Title)
	assert.Equal(t, "barbell", tmpl.Equipment)
	assert.False(t, tmpl.IsCustom)
}

func TestExerciseTemplatesGet(t *testing.T) {
	const id = "3BC06AD3"
	_, client := newTestServer(t, "/v1/exercise_templates/"+id, "exercise_template_get.json", 200)
	tmpl, err := client.ExerciseTemplates.Get(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, tmpl.ID)
	assert.Equal(t, "21s Bicep Curl", tmpl.Title)
	assert.Equal(t, "weight_reps", tmpl.Type)
	assert.Equal(t, "biceps", tmpl.PrimaryMuscleGroup)
	assert.Equal(t, "barbell", tmpl.Equipment)
	assert.False(t, tmpl.IsCustom)
}

func TestExerciseTemplatesInvalidPageSize(t *testing.T) {
	_, client := newTestServer(t, "/v1/exercise_templates", "exercise_templates_list.json", 200)
	_, err := client.ExerciseTemplates.List(context.Background(), 1, 101)
	assert.ErrorIs(t, err, hevy.ErrInvalidPageSize)
}

func TestExerciseLimitExceeded(t *testing.T) {
	_, client := newErrorServer(t, 403)
	_, err := client.ExerciseTemplates.Create(context.Background(), hevy.CreateExerciseInput{
		Title:             "Test",
		ExerciseType:      hevy.ExerciseTypeWeightReps,
		EquipmentCategory: hevy.EquipmentBarbell,
		MuscleGroup:       hevy.MuscleGroupChest,
	})
	assert.ErrorIs(t, err, hevy.ErrExerciseLimitExceeded)
}
