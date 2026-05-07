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

func TestRoutineFolderPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")

		file := fmt.Sprintf("testdata/responses/routine_folder/routine_folder-%s.json", page)
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		_, err = res.Write(data)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	routineFolders := []hevy.RoutineFolder{}
	pager := client.RoutineFolders()
	for x, err := range pager {
		assert.NoError(t, err)
		routineFolders = append(routineFolders, x)
	}

	assert.NotEmpty(t, routineFolders)

	assert.Len(t, routineFolders, 7)
}

func TestRoutineFolder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/routine_folders":
			var file string
			if req.Method == http.MethodGet {
				page := req.URL.Query().Get("page")
				file = fmt.Sprintf("testdata/responses/routine_folder/routine_folder-%s.json", page)
			} else {
				file = "testdata/responses/routine_folder/create_folder.json"
			}
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)

		case "/v1/routine_folders/1273009":
			data, err := os.ReadFile("testdata/responses/routine_folder/single-routine_folder.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	t.Run("Test All Routine Folders", func(t *testing.T) {
		routineFolders, err := client.AllRoutineFolders()
		assert.NoError(t, err)
		assert.NotEmpty(t, routineFolders)
		assert.Len(t, routineFolders, 7)
	})

	t.Run("Test Single Routine Folder", func(t *testing.T) {
		routineFolder, err := client.RoutineFolder(1273009)
		assert.NoError(t, err)
		assert.NotEmpty(t, routineFolder)
		assert.Equal(t, 1273009, routineFolder.ID)
		assert.Equal(t, 4, routineFolder.Index)
		assert.Equal(t, "Intermediate Full-Body (Gym Equipment)", routineFolder.Title)
	})

	t.Run("Test Get Routine Folders", func(t *testing.T) {
		routineFolders, next, err := client.GetRoutineFolders(2, 2)
		assert.Equal(t, 3, next)
		assert.NoError(t, err)
		assert.NotEmpty(t, routineFolders)
		assert.Len(t, routineFolders, 2)
	})

	t.Run("Test Create Routine Folder", func(t *testing.T) {
		routineFolder, err := client.CreateRoutineFolder("New Routine Folder")
		assert.NoError(t, err)
		assert.NotEmpty(t, routineFolder)
		assert.Equal(t, 2827406, routineFolder.ID)
		assert.Equal(t, 0, routineFolder.Index)
		assert.Equal(t, "Testing Folder", routineFolder.Title)
	})
}
