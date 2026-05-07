package hevy

import "fmt"

type exerciseTemplateResponse struct {
	paginatedResults
	ExerciseTemplates []ExerciseTemplate `json:"exercise_templates"`
}

// ExerciseTemplates returns an iterator that yields exercise templates one by one.
// If an error occurs fetching a page, it is yielded as the second value and iteration stops.
func (c Client) ExerciseTemplates() func(func(ExerciseTemplate, error) bool) {
	size := 10
	return func(yield func(ExerciseTemplate, error) bool) {
		page := 1

		for {
			resp, next, err := c.GetExerciseTemplates(page, size)
			if err != nil {
				yield(ExerciseTemplate{}, err)
				return
			}

			for _, template := range resp {
				if !yield(template, nil) {
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

// AllExerciseTemplates gets all exercise templates.
// This is a convenience method that handles pagination for you, if you have a large
// number of exercise templates this may take a while to complete.
func (c Client) AllExerciseTemplates() ([]ExerciseTemplate, error) {
	templates := []ExerciseTemplate{}

	page := 1
	size := 100

	for {
		resp, next, err := c.GetExerciseTemplates(page, size)
		if err != nil {
			return nil, err
		}

		templates = append(templates, resp...)
		if next == 0 {
			break
		}
		page = next
	}

	return templates, nil
}

// GetExerciseTemplates gets a paged list of exercise templates. Page is the paginated page
// number (starting at 1) and size is the number of exercise templates to return per page (max 100).
// It returns the list of exercise templates, the next page number (or 0 if there are no more pages)
// and an error if one occurred.
func (c Client) GetExerciseTemplates(page, size int) ([]ExerciseTemplate, int, error) {
	if size > 100 {
		size = 100
	}
	q := map[string]string{
		"page":     fmt.Sprintf("%d", page),
		"pageSize": fmt.Sprintf("%d", size),
	}
	url := c.constructURL("exercise_templates", q)
	result := exerciseTemplateResponse{}
	err := c.get(url, &result)
	if err != nil {
		return nil, 0, err
	}

	next := result.Page + 1
	if result.Page >= result.PageCount {
		next = 0
	}
	return result.ExerciseTemplates, next, nil
}

// ExerciseTemplate returns an exercise template by ID.
func (c Client) ExerciseTemplate(id string) (ExerciseTemplate, error) {
	url := c.constructURL(fmt.Sprintf("exercise_templates/%s", id), nil)

	result := ExerciseTemplate{}
	err := c.get(url, &result)
	if err != nil {
		return ExerciseTemplate{}, err
	}
	return result, nil
}
