package hevy

import "fmt"

type routineFolderResponse struct {
	paginatedResults
	RoutineFolders []RoutineFolder `json:"routine_folders"`
}

// RoutineFolders returns an iterator that yields routine folders one by one.
func (c Client) RoutineFolders() func(func(RoutineFolder) bool) {
	size := 10
	return func(yield func(RoutineFolder) bool) {
		page := 1

		for {
			resp, next, err := c.GetRoutineFolders(page, size)
			if err != nil {
				return
			}

			for _, routineFolder := range resp {
				if !yield(routineFolder) {
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

// AllRoutineFolders gets all routine folders. This is a convenience method that handles pagination
// for you, if you have a large number of routine folders this may take a while to complete.
func (c Client) AllRoutineFolders() ([]RoutineFolder, error) {
	routineFolders := []RoutineFolder{}

	page := 1
	size := 10
	for {
		resp, next, err := c.GetRoutineFolders(page, size)
		if err != nil {
			return nil, err
		}

		routineFolders = append(routineFolders, resp...)

		if next == 0 {
			break
		}
		page = next
	}

	return routineFolders, nil
}

// GetRoutineFolders gets a page of routine folders. The next page can be fetched by using the returned
// next value, if next is 0 there are no more pages to fetch.
func (c Client) GetRoutineFolders(page int, size int) ([]RoutineFolder, int, error) {
	if size > 10 {
		size = 10
	}

	q := map[string]string{
		"page":     fmt.Sprintf("%d", page),
		"pageSize": fmt.Sprintf("%d", size),
	}

	url := c.constructURL("routine_folders", q)
	result := routineFolderResponse{}
	err := c.get(url, &result)
	if err != nil {
		return nil, 0, err
	}

	next := result.Page + 1
	if result.Page >= result.PageCount {
		next = 0
	}
	return result.RoutineFolders, next, nil
}

// RoutineFolder gets a routine folder by id.
func (c Client) RoutineFolder(id int) (RoutineFolder, error) {
	path := fmt.Sprintf("routine_folders/%d", id)
	url := c.constructURL(path, map[string]string{})
	result := RoutineFolder{}
	err := c.get(url, &result)
	if err != nil {
		return RoutineFolder{}, err
	}

	return result, nil
}
