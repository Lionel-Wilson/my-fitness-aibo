package mappers

type SimpleErrorResponse struct {
	Error string `json:"error"`
}

type SimpleMessageResponse struct {
	Message string `json:"message"`
}

func ToSimpleErrorResponse(msg string) SimpleErrorResponse {
	return SimpleErrorResponse{Error: msg}
}

func ToSimpleMessageResponse(msg string) SimpleMessageResponse {
	return SimpleMessageResponse{Message: msg}
}
