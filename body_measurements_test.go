package hevy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestBodyMeasurementsList(t *testing.T) {
	_, client := newTestServer(t, "/v1/body_measurements", "body_measurements_list.json", 200)
	page, err := client.BodyMeasurements.List(context.Background(), 1, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, page.PageCount)
	require.Len(t, page.Measurements, 2)

	m := page.Measurements[0]
	assert.Equal(t, "2026-02-26", m.Date)
	require.NotNil(t, m.WeightKg)
	assert.Equal(t, 82.19103518973792, *m.WeightKg)
}

func TestBodyMeasurementsGet(t *testing.T) {
	_, client := newTestServer(t, "/v1/body_measurements/2026-02-26", "body_measurement_get.json", 200)
	m, err := client.BodyMeasurements.Get(context.Background(), "2026-02-26")
	require.NoError(t, err)
	assert.Equal(t, "2026-02-26", m.Date)
	require.NotNil(t, m.WeightKg)
	assert.Equal(t, 82.19103518973792, *m.WeightKg)
}

func TestBodyMeasurementsConflict(t *testing.T) {
	_, client := newErrorServer(t, 409)
	err := client.BodyMeasurements.Create(context.Background(), hevy.BodyMeasurement{
		Date: "2026-02-26",
	})
	assert.ErrorIs(t, err, hevy.ErrConflict)
}

func TestBodyMeasurementsGetNotFound(t *testing.T) {
	_, client := newErrorServer(t, 404)
	_, err := client.BodyMeasurements.Get(context.Background(), "2000-01-01")
	assert.ErrorIs(t, err, hevy.ErrNotFound)
}

func TestBodyMeasurementsListInvalidPageSize(t *testing.T) {
	_, client := newTestServer(t, "/v1/body_measurements", "body_measurements_list.json", 200)
	_, err := client.BodyMeasurements.List(context.Background(), 1, 11)
	assert.ErrorIs(t, err, hevy.ErrInvalidPageSize)
}
