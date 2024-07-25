package entities

type Response struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
