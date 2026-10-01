package handler

type APIError struct {
	Status  int               `json:"status"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type Response[T any] struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Data    T         `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}
