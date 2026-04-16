package res

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
