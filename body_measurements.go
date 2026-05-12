package hevy

import (
	"context"
	"errors"
)

// BodyMeasurementsService handles body measurement endpoints.
type BodyMeasurementsService struct{ c *core }

// List returns one page of body measurements. page is 1-based; pageSize max is 10.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *BodyMeasurementsService) List(ctx context.Context, page, pageSize int) (*PaginatedBodyMeasurements, error) {
	if err := validatePageSize(pageSize, 10); err != nil {
		return nil, err
	}
	var out PaginatedBodyMeasurements
	if err := s.c.get(ctx, "/v1/body_measurements", pageQuery(page, pageSize), &out); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNoMorePages
		}
		return nil, err
	}
	if page > out.PageCount && out.PageCount > 0 {
		return nil, ErrNoMorePages
	}
	return &out, nil
}

// ListAll fetches every page and returns all body measurements in a single slice.
func (s *BodyMeasurementsService) ListAll(ctx context.Context) ([]BodyMeasurement, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]BodyMeasurement, int, error) {
		p, err := s.List(ctx, page, 10)
		if err != nil {
			return nil, 0, err
		}
		return p.Measurements, p.PageCount, nil
	})
}

// Get returns the body measurement entry for a specific date (format: "YYYY-MM-DD").
func (s *BodyMeasurementsService) Get(ctx context.Context, date string) (*BodyMeasurement, error) {
	var out BodyMeasurement
	if err := s.c.get(ctx, "/v1/body_measurements/"+date, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new body measurement entry.
// Returns ErrConflict if an entry already exists for that date.
func (s *BodyMeasurementsService) Create(ctx context.Context, m BodyMeasurement) error {
	return s.c.post(ctx, "/v1/body_measurements", m, nil)
}

// Update replaces all fields for the measurement on a given date (format: "YYYY-MM-DD").
// Fields set to nil will be cleared to null on the server.
func (s *BodyMeasurementsService) Update(ctx context.Context, date string, m BodyMeasurementUpdate) error {
	return s.c.put(ctx, "/v1/body_measurements/"+date, m, nil)
}
