package hevy_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/user":
			data, err := os.ReadFile("testdata/responses/user/info.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	user, err := client.User()
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, "d2670e5a-951c-4911-9d12-43e6eff605b1", user.ID.String())
	assert.Equal(t, "khabiaz", user.Name)
	assert.Equal(t, "https://hevy.com/user/khabiaz", user.URL)
}
