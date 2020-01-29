package controller

type UriIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

// swagger:model
type CountResp struct {
	// example: 1
	Count int `json:"count"`
}
