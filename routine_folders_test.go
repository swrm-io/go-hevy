package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestRoutineFoldersList(t *testing.T) {
	_, client := newTestServer(t, "/v1/routine_folders", "routine_folders_list.json", 200)
	page, err := client.RoutineFolders.List(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 4, page.PageCount)
	require.Len(t, page.RoutineFolders, 2)
	assert.Equal(t, 2014654, page.RoutineFolders[0].ID)
	assert.Equal(t, "AI", page.RoutineFolders[0].Title)
}

func TestRoutineFoldersGet(t *testing.T) {
	_, client := newTestServer(t, "/v1/routine_folders/2014654", "routine_folder_get.json", 200)
	f, err := client.RoutineFolders.Get(context.Background(), 2014654)
	require.NoError(t, err)
	assert.Equal(t, 2014654, f.ID)
	assert.Equal(t, "AI", f.Title)
	assert.Equal(t, 15, f.Index)
}

func TestRoutineFoldersListInvalidPageSize(t *testing.T) {
	_, client := newTestServer(t, "/v1/routine_folders", "routine_folders_list.json", 200)
	_, err := client.RoutineFolders.List(context.Background(), 1, 11)
	assert.ErrorIs(t, err, hevy.ErrInvalidPageSize)
}
