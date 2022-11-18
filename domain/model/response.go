package model

type Response struct {
	Payload    interface{} `json:"payload"`
	TrackingId string      `json:"trackingId"`
	Error      string      `json:"error"`
	Success    bool        `json:"success"`
}
