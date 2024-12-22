package entities

type JobDTO struct {
	Code    int
	Status  string
	Message string
}

type MessengerClientResponse struct {
	Status       int `json:"status"`
	InternalCode int `json:"internalCode"`
	Data         []struct {
		ID      string `json:"id"`
		To      string `json:"to"`
		Content string `json:"content"`
	} `json:"data"`
}
