package models

// swagger:model
type TaskFuzzResult struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`
	// example: /afl/afl-fuzz -i xx -o xx ./test
	Command string `json:"command"`
	// example: 60
	TimeExecuted int `json:"timeExecuted"`
	// example: 1
	TaskID uint64 `json:"taskid" sql:"type:integer REFERENCES task(id) ON DELETE CASCADE"`
	// example: 1579996805
	UpdateAt int64 `json:"updateAt"`
}

type TaskFuzzResultStat struct {
	Key              string `json:"key"`
	Value            string `json:"value"`
	TaskFuzzResultID uint64 `json:"taskid" sql:"type:integer REFERENCES task_fuzz_result(id) ON DELETE CASCADE"`
}
