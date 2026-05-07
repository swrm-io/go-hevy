package hevy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestExerciseTemplatePagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")

		file := fmt.Sprintf("testdata/responses/exercise_template/exercise_template-%s.json", page)
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		_, err = res.Write(data)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	templates := []hevy.ExerciseTemplate{}
	pager := client.ExerciseTemplates()
	for x := range pager {
		templates = append(templates, x)
	}

	assert.NotEmpty(t, templates)

	assert.Len(t, templates, 9)
}

func TestExerciseTemplate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/exercise_templates":
			page := req.URL.Query().Get("page")

			file := fmt.Sprintf("testdata/responses/exercise_template/exercise_template-%s.json", page)
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/exercise_templates/4F5866F8":
			data, err := os.ReadFile("testdata/responses/exercise_template/single-exercise_template.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	t.Run("Test All Exercise Templates", func(t *testing.T) {
		templates, err := client.AllExerciseTemplates()
		assert.NoError(t, err)
		assert.NotEmpty(t, templates)
		assert.Len(t, templates, 9)
	})

	t.Run("Test Get Exercise Template", func(t *testing.T) {
		template, err := client.ExerciseTemplate("4F5866F8")
		assert.NoError(t, err)
		assert.Equal(t, "4F5866F8", template.ID)
	})
}
