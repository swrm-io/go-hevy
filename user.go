package hevy

type infoResponse struct {
	Data User `json:"data"`
}

func (c Client) User() (User, error) {
	url := c.constructURL("user", map[string]string{})
	result := infoResponse{}
	err := c.get(url, &result)
	return result.Data, err
}
