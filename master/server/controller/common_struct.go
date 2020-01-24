package controller

type UriIDReq struct {
	ID uint64 `uri:"id" binding:"required"`
}
