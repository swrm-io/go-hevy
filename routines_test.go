package hevy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestRoutine(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/routines":
			page := req.URL.Query().Get("page")

			file := fmt.Sprintf("testdata/responses/routine/routine-%s.json", page)
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/routines/9c3e6a25-67e9-4c0c-860c-286e13fe9924":
			data, err := os.ReadFile("testdata/responses/routine/single-routine.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	t.Run("Test All Routines", func(t *testing.T) {
		routines, err := client.AllRoutines()
		assert.NoError(t, err)
		assert.NotEmpty(t, routines)
		assert.Len(t, routines, 3)
	})

	t.Run("Test Paginated Routines", func(t *testing.T) {
		routines, next, err := client.GetRoutines(2, 1)
		assert.Equal(t, 3, next)
		assert.NoError(t, err)
		assert.NotEmpty(t, routines)
		assert.Len(t, routines, 1)
	})

	t.Run("Test Single Routine", func(t *testing.T) {
		routineID, err := uuid.Parse("9c3e6a25-67e9-4c0c-860c-286e13fe9924")
		assert.NoError(t, err)

		routine, err := client.Routine(routineID)
		assert.NoError(t, err)
		assert.NotEmpty(t, routine)
		assert.Equal(t, "9c3e6a25-67e9-4c0c-860c-286e13fe9924", routine.ID.String())
	})
}
