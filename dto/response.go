package dto

type BaseResponse struct {
	Code int                    `json:"Code"`
	Msg  string                 `json:"Msg"`
	Data map[string]interface{} `json:"Data"`
}
