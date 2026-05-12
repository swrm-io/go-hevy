package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserInfo(t *testing.T) {
	client := newTestServer(t, "/v1/user/info", "user_info.json")
	user, err := client.User.Info(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "khabiaz", user.Name)
	assert.Equal(t, "d2670e5a-951c-4911-9d12-43e6eff605b1", user.ID)
	assert.Equal(t, "https://hevy.com/user/khabiaz", user.URL)
}
