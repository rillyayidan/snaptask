package model

type ExtractedItem struct {
	Type     string  `json:"type"`
	Title    string  `json:"title"`
	Detail   string  `json:"detail"`
	DueDate  *string `json:"due_date"`
	Priority string  `json:"priority"`
}

type ExtractResponse struct {
	Items []ExtractedItem `json:"items"`
}

type PushRequest struct {
	AccessToken string          `json:"access_token"`
	Items       []ExtractedItem `json:"items"`
}

type PushResult struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status string `json:"status"`
	ID     string `json:"id,omitempty"`
	Error  string `json:"error,omitempty"`
}

type PushResponse struct {
	Results []PushResult `json:"results"`
}
