package hevy

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// A response for fetching a list of workouts
type workoutResponse struct {
	paginatedResults
	Workouts []Workout `json:"workouts"`
}

type workoutCountResponse struct {
	Count int `json:"workout_count"`
}

type workoutEventResponse struct {
	paginatedResults
	Events []Event
}

// Workouts returns an iterator that yields workouts one by one.
func (c Client) Workouts() func(func(Workout) bool) {
	size := 10
	return func(yield func(Workout) bool) {
		page := 1

		for {
			resp, next, err := c.GetWorkouts(page, size)
			if err != nil {
				return
			}

			for _, workout := range resp {
				if !yield(workout) {
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

// AllWorkouts gets all workouts.
// This is a convenience method that handles pagination for you, if you have a large
// number of workouts this may take a while to complete.
func (c Client) AllWorkouts() ([]Workout, error) {
	workouts := []Workout{}

	page := 1
	size := 10

	for {
		resp, next, err := c.GetWorkouts(page, size)
		if err != nil {
			return nil, err
		}

		workouts = append(workouts, resp...)
		if next == 0 {
			break
		}
		page = next
	}

	return workouts, nil
}

// GetWorkouts retrieves a paged list of workouts. Page is the paginated page number (starting at 1) and size is the number of
// workouts to return per page. The maximum page size is 10, limited by the API.
// It returns the list of workouts, the next page (0 if there are no more pages) and an error if one occurred.
// Workouts are returned ordered from newest to oldest.
func (c Client) GetWorkouts(page int, size int) ([]Workout, int, error) {
	if size > 10 {
		size = 10
	}
	q := map[string]string{
		"page":     fmt.Sprintf("%d", page),
		"pageSize": fmt.Sprintf("%d", size),
	}
	url := c.constructURL("workouts", q)
	result := workoutResponse{}
	err := c.get(url, &result)
	if err != nil {
		return nil, 0, err
	}

	next := result.Page + 1
	if result.Page >= result.PageCount {
		next = 0
	}
	return result.Workouts, next, nil
}

// Workout retrieves a single workout by its ID.
func (c Client) Workout(id uuid.UUID) (Workout, error) {
	path := fmt.Sprintf("workouts/%s", id.String())
	url := c.constructURL(path, map[string]string{})

	result := Workout{}

	err := c.get(url, &result)
	if err != nil {
		return Workout{}, err
	}

	return result, nil
}

// WorkoutCount returns a count of workouts
func (c Client) WorkoutCount() (int, error) {
	url := c.constructURL("workouts/count", map[string]string{})

	result := workoutCountResponse{}

	err := c.get(url, &result)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}

// WorkoutEvents retrieves a paged list of workout events (updates or deletes) since a given date.
// Events are ordered from newest to oldest. The intention is to allow clients to keep their local
// cache of workouts up to date without having to fetch the entire list of workouts.
func (c Client) WorkoutEvents(since time.Time) ([]Event, error) {
	events := []Event{}

	page := 1
	size := 10

	for {
		q := map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"pageSize": fmt.Sprintf("%d", size),
			"since":    since.Format("RFC3339Nano"),
		}
		url := c.constructURL("workouts/events", q)
		result := workoutEventResponse{}
		err := c.get(url, &result)
		if err != nil {
			return nil, err
		}

		events = append(events, result.Events...)

		if result.Page == result.PageCount {
			break
		}
		page++
	}

	return events, nil
}
