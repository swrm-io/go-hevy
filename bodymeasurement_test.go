package hevy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/swrm-io/go-hevy"
)

func TestBodyMeasurementPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")

		file := fmt.Sprintf("testdata/responses/body_measurement/body_measurement-%s.json", page)
		data, err := os.ReadFile(file)

		assert.NoError(t, err)
		_, err = res.Write(data)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	bodyMeasurements := []hevy.BodyMeasurement{}
	pager := client.BodyMeasurements()
	for x := range pager {
		bodyMeasurements = append(bodyMeasurements, x)
	}

	assert.NotEmpty(t, bodyMeasurements)

	assert.Len(t, bodyMeasurements, 5)
}

func TestBodyMeasurement(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/body_measurements":
			page := req.URL.Query().Get("page")

			file := fmt.Sprintf("testdata/responses/body_measurement/body_measurement-%s.json", page)
			data, err := os.ReadFile(file)
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		case "/v1/body_measurements/2024-11-04":
			data, err := os.ReadFile("testdata/responses/body_measurement/single-body_measurement.json")
			assert.NoError(t, err)
			_, err = res.Write(data)
			assert.NoError(t, err)
		}
	}))
	defer srv.Close()

	client := hevy.NewClient("my-fake-api-key")
	client.APIURL = srv.URL

	t.Run("Test All Body Measurements", func(t *testing.T) {
		bodyMeasurements, err := client.AllBodyMeasurements()
		assert.NoError(t, err)
		assert.NotEmpty(t, bodyMeasurements)
		assert.Len(t, bodyMeasurements, 5)
	})

	t.Run("Test Get Body Measurement", func(t *testing.T) {
		bodyMeasurement, next, err := client.GetBodyMeasurements(2, 2)
		assert.Equal(t, 3, next)
		assert.NoError(t, err)
		assert.NotEmpty(t, bodyMeasurement)
		assert.Len(t, bodyMeasurement, 2)
	})

	t.Run("Test Single Body Measurement", func(t *testing.T) {
		bodyMeasurement, err := client.BodyMeasurement(time.Date(2024, 11, 4, 0, 0, 0, 0, time.UTC))
		assert.NoError(t, err)
		assert.NotEmpty(t, bodyMeasurement)
		assert.Equal(t, 6648419, bodyMeasurement.ID)
		assert.Equal(t, "2024-11-04", bodyMeasurement.Date)
		assert.Equal(t, time.Date(2024, 11, 5, 6, 37, 16, 243000000, time.UTC), bodyMeasurement.CreatedAt)
		assert.Equal(t, 82.19103518973792, bodyMeasurement.WeightKG)
		assert.Equal(t, "181.20", fmt.Sprintf("%.2f", bodyMeasurement.WeightLB))
		assert.Equal(t, "38.89", fmt.Sprintf("%.2f", bodyMeasurement.ChestIn))
	})
}
