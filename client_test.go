package hevy_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestAPIErrorDetails(t *testing.T) {
	client := newErrorServer(t, 400)
	_, err := client.Workouts.Get(context.Background(), "bad-id")
	require.Error(t, err)

	var apiErr *hevy.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.ErrorIs(t, err, hevy.ErrBadRequest)
}

func TestAPIErrorUnwrap(t *testing.T) {
	client := newErrorServer(t, 401)
	_, err := client.User.Info(context.Background())
	assert.ErrorIs(t, err, hevy.ErrUnauthorized)

	var apiErr *hevy.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 401, apiErr.StatusCode)
	assert.True(t, errors.Is(apiErr, hevy.ErrUnauthorized))
}
