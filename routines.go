package hevy

import (
	"context"
	"errors"
)

// RoutinesService handles routine endpoints.
type RoutinesService struct{ c *core }

// List returns one page of routines. page is 1-based; pageSize max is 10.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *RoutinesService) List(ctx context.Context, page, pageSize int) (*PaginatedRoutines, error) {
	if err := validatePageSize(pageSize, 10); err != nil {
		return nil, err
	}
	var out PaginatedRoutines
	if err := s.c.get(ctx, "/v1/routines", pageQuery(page, pageSize), &out); err != nil {
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

// ListAll fetches every page and returns all routines in a single slice.
func (s *RoutinesService) ListAll(ctx context.Context) ([]Routine, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]Routine, int, error) {
		p, err := s.List(ctx, page, 10)
		if err != nil {
			return nil, 0, err
		}
		return p.Routines, p.PageCount, nil
	})
}

// Get returns a single routine by ID.
func (s *RoutinesService) Get(ctx context.Context, routineID string) (*Routine, error) {
	var out struct {
		Routine Routine `json:"routine"`
	}
	if err := s.c.get(ctx, "/v1/routines/"+routineID, nil, &out); err != nil {
		return nil, err
	}
	return &out.Routine, nil
}

// Create creates a new routine and returns it.
// Returns ErrRoutineLimitExceeded if the account routine limit is reached.
func (s *RoutinesService) Create(ctx context.Context, routine RoutineInput) (*Routine, error) {
	var out Routine
	if err := s.c.post(ctx, "/v1/routines", map[string]any{"routine": routine}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates an existing routine and returns the updated version.
func (s *RoutinesService) Update(ctx context.Context, routineID string, routine RoutineUpdateInput) (*Routine, error) {
	var out Routine
	if err := s.c.put(ctx, "/v1/routines/"+routineID, map[string]any{"routine": routine}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
