package dto

type SetConfigRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type GetConfigResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
