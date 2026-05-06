package hevy

import (
	"fmt"

	"github.com/google/uuid"
)

// A response for fetching a list of workouts
type routineResponse struct {
	paginatedResults
	Routines []Routine `json:"routines"`
}

func (c Client) Routines() func(func(Routine) bool) {
	size := 10
	return func(yield func(Routine) bool) {
		page := 1

		for {
			resp, next, err := c.GetRoutines(page, size)
			if err != nil {
				return
			}

			for _, routine := range resp {
				if !yield(routine) {
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

func (c Client) AllRoutines() ([]Routine, error) {
	routines := []Routine{}

	page := 1
	size := 10

	for {
		resp, next, err := c.GetRoutines(page, size)
		if err != nil {
			return nil, err
		}

		routines = append(routines, resp...)

		if next == 0 {
			break
		}
		page = next
	}

	return routines, nil
}

// Routines gets all routines.
func (c Client) GetRoutines(page int, size int) ([]Routine, int, error) {
	if size > 10 {
		size = 10
	}

	q := map[string]string{
		"page":     fmt.Sprintf("%d", page),
		"pageSize": fmt.Sprintf("%d", size),
	}
	url := c.constructURL("routines", q)
	result := routineResponse{}
	err := c.get(url, &result)
	if err != nil {
		return nil, 0, err
	}

	next := result.Page + 1
	if result.Page == result.PageCount {
		next = 0
	}

	return result.Routines, next, nil
}

func (c Client) Routine(id uuid.UUID) (Routine, error) {
	path := fmt.Sprintf("routines/%s", id.String())
	url := c.constructURL(path, map[string]string{})

	result := Routine{}

	err := c.get(url, &result)
	if err != nil {
		return Routine{}, err
	}

	return result, nil
}
