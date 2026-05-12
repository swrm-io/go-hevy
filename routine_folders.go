package hevy

import (
	"context"
	"errors"
	"fmt"
)

// RoutineFoldersService handles routine folder endpoints.
type RoutineFoldersService struct{ c *core }

// List returns one page of routine folders. page is 1-based; pageSize max is 10.
// Returns ErrNoMorePages when page exceeds the total page count.
func (s *RoutineFoldersService) List(ctx context.Context, page, pageSize int) (*PaginatedRoutineFolders, error) {
	if err := validatePageSize(pageSize, 10); err != nil {
		return nil, err
	}
	var out PaginatedRoutineFolders
	if err := s.c.get(ctx, "/v1/routine_folders", pageQuery(page, pageSize), &out); err != nil {
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

// ListAll fetches every page and returns all routine folders in a single slice.
func (s *RoutineFoldersService) ListAll(ctx context.Context) ([]RoutineFolder, error) {
	return listAll(ctx, func(ctx context.Context, page int) ([]RoutineFolder, int, error) {
		p, err := s.List(ctx, page, 10)
		if err != nil {
			return nil, 0, err
		}
		return p.RoutineFolders, p.PageCount, nil
	})
}

// Get returns a single routine folder by ID.
func (s *RoutineFoldersService) Get(ctx context.Context, folderID int) (*RoutineFolder, error) {
	var out RoutineFolder
	if err := s.c.get(ctx, fmt.Sprintf("/v1/routine_folders/%d", folderID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new routine folder with the given title.
func (s *RoutineFoldersService) Create(ctx context.Context, title string) (*RoutineFolder, error) {
	body := map[string]any{"routine_folder": map[string]string{"title": title}}
	var out RoutineFolder
	if err := s.c.post(ctx, "/v1/routine_folders", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
