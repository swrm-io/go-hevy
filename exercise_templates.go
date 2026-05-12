package hevy

import (
	"context"
	"errors"
)

// ExerciseTemplatesService handles exercise template endpoints.
type ExerciseTemplatesService struct{ c *core }

// List returns one page of exercise templates. page is 1-based; pageSize max is 100.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *ExerciseTemplatesService) List(ctx context.Context, page, pageSize int) (*PaginatedExerciseTemplates, error) {
	if err := validatePageSize(pageSize, 100); err != nil {
		return nil, err
	}
	var out PaginatedExerciseTemplates
	if err := s.c.get(ctx, "/v1/exercise_templates", pageQuery(page, pageSize), &out); err != nil {
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

// ListAll fetches every page and returns all exercise templates in a single slice.
func (s *ExerciseTemplatesService) ListAll(ctx context.Context) ([]ExerciseTemplate, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]ExerciseTemplate, int, error) {
		p, err := s.List(ctx, page, 100)
		if err != nil {
			return nil, 0, err
		}
		return p.ExerciseTemplates, p.PageCount, nil
	})
}

// Get returns a single exercise template by ID.
func (s *ExerciseTemplatesService) Get(ctx context.Context, exerciseTemplateID string) (*ExerciseTemplate, error) {
	var out ExerciseTemplate
	if err := s.c.get(ctx, "/v1/exercise_templates/"+exerciseTemplateID, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a custom exercise template and returns its ID.
// Returns ErrExerciseLimitExceeded if the account limit is reached.
func (s *ExerciseTemplatesService) Create(ctx context.Context, exercise CreateExerciseInput) (int, error) {
	var out struct {
		ID int `json:"id"`
	}
	if err := s.c.post(ctx, "/v1/exercise_templates", map[string]any{"exercise": exercise}, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
