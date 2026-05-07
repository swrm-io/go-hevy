package hevy

import (
	"fmt"
	"time"
)

type bodyMeasurementsResponse struct {
	paginatedResults
	BodyMeasurements []BodyMeasurement `json:"body_measurements"`
}

// BodyMeasurements returns an iterator that yields body measurements one by one.
// If an error occurs fetching a page, it is yielded as the second value and iteration stops.
func (c Client) BodyMeasurements() func(func(BodyMeasurement, error) bool) {
	size := 10
	return func(yield func(BodyMeasurement, error) bool) {
		page := 1

		for {
			resp, next, err := c.GetBodyMeasurements(page, size)
			if err != nil {
				yield(BodyMeasurement{}, err)
				return
			}

			for _, measurement := range resp {
				if !yield(measurement, nil) {
					return
				}
			}

			if next == 0 {
				break
			}
			page++
		}
	}
}

func (c Client) AllBodyMeasurements() ([]BodyMeasurement, error) {
	measurements := []BodyMeasurement{}

	page := 1
	size := 10

	for {
		resp, next, err := c.GetBodyMeasurements(page, size)
		if err != nil {
			return nil, err
		}

		measurements = append(measurements, resp...)

		if next == 0 {
			break
		}
		page = next
	}

	return measurements, nil
}

func (c Client) GetBodyMeasurements(page int, size int) ([]BodyMeasurement, int, error) {
	if size > 10 {
		size = 10
	}

	q := map[string]string{
		"page":     fmt.Sprintf("%d", page),
		"pageSize": fmt.Sprintf("%d", size),
	}

	url := c.constructURL("body_measurements", q)
	result := bodyMeasurementsResponse{}
	err := c.get(url, &result)
	if err != nil {
		return nil, 0, err
	}

	next := result.Page + 1
	if result.Page >= result.PageCount {
		next = 0
	}
	return result.BodyMeasurements, next, nil
}

func (c Client) BodyMeasurement(date time.Time) (BodyMeasurement, error) {
	path := fmt.Sprintf("body_measurements/%s", date.Format("2006-01-02"))
	url := c.constructURL(path, map[string]string{})

	result := BodyMeasurement{}

	err := c.get(url, &result)
	if err != nil {
		return BodyMeasurement{}, err
	}

	return result, nil
}
