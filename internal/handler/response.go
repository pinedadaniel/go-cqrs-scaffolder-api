package handler

func Success[T any](data T) Response[T] {
	return Response[T]{
		Success: true,
		Message: "Request Successful",
		Data:    data,
	}
}

func Error(code, message string, details map[string]string) Response[any] {
	return Response[any]{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
