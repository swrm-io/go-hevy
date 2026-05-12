package hevy

import (
	"context"
	"net/url"
	"time"
)

// ExerciseHistoryService handles exercise history endpoints.
type ExerciseHistoryService struct{ c *core }

// GetOptions are optional date range filters for exercise history.
type GetHistoryOptions struct {
	// StartDate filters entries to those on or after this date.
	StartDate *time.Time
	// EndDate filters entries to those on or before this date.
	EndDate *time.Time
}

// Get returns the full history of sets performed for a given exercise template.
func (s *ExerciseHistoryService) Get(ctx context.Context, exerciseTemplateID string, opts *GetHistoryOptions) ([]ExerciseHistoryEntry, error) {
	q := url.Values{}
	if opts != nil {
		if opts.StartDate != nil {
			q.Set("start_date", opts.StartDate.UTC().Format(time.RFC3339))
		}
		if opts.EndDate != nil {
			q.Set("end_date", opts.EndDate.UTC().Format(time.RFC3339))
		}
	}
	var out ExerciseHistory
	if err := s.c.get(ctx, "/v1/exercise_history/"+exerciseTemplateID, q, &out); err != nil {
		return nil, err
	}
	return out.ExerciseHistory, nil
}
