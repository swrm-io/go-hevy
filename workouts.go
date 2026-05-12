package hevy

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// WorkoutsService handles workout endpoints.
type WorkoutsService struct{ c *core }

// List returns one page of workouts. page is 1-based; pageSize max is 10.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *WorkoutsService) List(ctx context.Context, page, pageSize int) (*PaginatedWorkouts, error) {
	if err := validatePageSize(pageSize, 10); err != nil {
		return nil, err
	}
	var out PaginatedWorkouts
	if err := s.c.get(ctx, "/v1/workouts", pageQuery(page, pageSize), &out); err != nil {
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

// ListAll fetches every page and returns all workouts in a single slice.
func (s *WorkoutsService) ListAll(ctx context.Context) ([]Workout, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]Workout, int, error) {
		p, err := s.List(ctx, page, 10)
		if err != nil {
			return nil, 0, err
		}
		return p.Workouts, p.PageCount, nil
	})
}

// Count returns the total number of workouts for the authenticated user.
func (s *WorkoutsService) Count(ctx context.Context) (int, error) {
	var out WorkoutCount
	if err := s.c.get(ctx, "/v1/workouts/count", nil, &out); err != nil {
		return 0, err
	}
	return out.WorkoutCount, nil
}

// WorkoutEventsOptions are optional filters for Events.
type WorkoutEventsOptions struct {
	// Since filters events to those occurring after this time.
	Since *time.Time
}

// Events returns one page of workout update/delete events. pageSize max is 10.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *WorkoutsService) Events(ctx context.Context, page, pageSize int, opts *WorkoutEventsOptions) (*PaginatedWorkoutEvents, error) {
	if err := validatePageSize(pageSize, 10); err != nil {
		return nil, err
	}
	q := pageQuery(page, pageSize)
	if opts != nil && opts.Since != nil {
		q.Set("since", opts.Since.UTC().Format(time.RFC3339))
	}
	var out PaginatedWorkoutEvents
	if err := s.c.get(ctx, "/v1/workouts/events", q, &out); err != nil {
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

// Get returns a single workout by ID.
func (s *WorkoutsService) Get(ctx context.Context, workoutID string) (*Workout, error) {
	var out Workout
	if err := s.c.get(ctx, "/v1/workouts/"+workoutID, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new workout and returns it.
func (s *WorkoutsService) Create(ctx context.Context, workout WorkoutInput) (*Workout, error) {
	var out Workout
	if err := s.c.post(ctx, "/v1/workouts", map[string]any{"workout": workout}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates an existing workout and returns the updated version.
func (s *WorkoutsService) Update(ctx context.Context, workoutID string, workout WorkoutInput) (*Workout, error) {
	var out Workout
	if err := s.c.put(ctx, "/v1/workouts/"+workoutID, map[string]any{"workout": workout}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// listAll is a generic helper that iterates all pages of a paginated resource.
func listAll[T any](ctx context.Context, fetch func(ctx context.Context, page int) ([]T, int, error)) ([]T, error) {
	var all []T
	page := 1
	for {
		items, pageCount, err := fetch(ctx, page)
		if errors.Is(err, ErrNoMorePages) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", page, err)
		}
		all = append(all, items...)
		if page >= pageCount {
			break
		}
		page++
	}
	return all, nil
}

// eventsListAll fetches all pages of workout events.
func (s *WorkoutsService) EventsAll(ctx context.Context, opts *WorkoutEventsOptions) ([]WorkoutEvent, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]WorkoutEvent, int, error) {
		p, err := s.Events(ctx, page, 10, opts)
		if err != nil {
			return nil, 0, err
		}
		return p.Events, p.PageCount, nil
	})
}

