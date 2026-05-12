package hevy

import "context"

// UserService handles user info endpoints.
type UserService struct{ c *core }

// Info returns basic information about the authenticated user.
func (s *UserService) Info(ctx context.Context) (*UserInfo, error) {
	var out struct {
		Data UserInfo `json:"data"`
	}
	if err := s.c.get(ctx, "/v1/user/info", nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}
