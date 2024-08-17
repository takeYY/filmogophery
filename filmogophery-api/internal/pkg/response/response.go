package response

type (
	OK struct {
		Message string `json:"message"`
	}
	ErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)
