package kubernetes

type clientStorageItemPostResp struct {
	ID            uint64 `json:"id" binding:"required"`
	Type          string `json:"type" binding:"required"`
	ExistsInImage bool   `json:"existsInImage" binding:"required"`
}

type clientTaskGetResp struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"errorMsg"`
}

type clientCrashGetResp struct {
	ID            uint64 `json:"id" binding:"required"`
	FileName      string `json:"fileName" binding:"required"`
	ReproduceAble bool   `json:"reproduceAble" binding:"required"`
}

type clientResultGetResp struct {
	Command      string            `json:"command" binding:"required"`
	TimeExecuted int               `json:"timeExecuted" binding:"required"`
	UpdateAt     int64             `json:"updateAt" binding:"required"`
	Stats        map[string]string `json:"stats" binding:"required"`
}
