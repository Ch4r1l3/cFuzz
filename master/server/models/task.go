package models

// swagger:model
type Task struct {
	// example: 1
	ID uint64 `gorm:"primary_key" json:"id"`

	// example: test
	Name string `json:"name" sql:"type:varchar(255) NOT NULL UNIQUE"`

	// example: 1
	ImageID uint64 `json:"imageID" sql:"type:integer REFERENCES image(id)"`

	// example: 60
	Time uint64 `json:"time"`

	// example: 60
	FuzzCycleTime uint64 `json:"fuzzCycleTime"`

	// example: 1
	FuzzerID uint64 `json:"fuzzerID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: 2
	CorpusID uint64 `json:"corpusID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: 3
	TargetID uint64 `json:"targetID" sql:"type:integer REFERENCES storage_item(id)"`

	// example: TaskRunning
	Status string `json:"status"`

	// example: pull image error
	ErrorMsg string `json:"errorMsg"`

	// example: 1579996805
	StatusUpdateAt int64 `json:"-"`

	// example: 1579996805
	StartedAt int64 `json:"startedAt"`

	// example: http://127.0.0.1/callback
	CallbackUrl string `json:"callbackUrl"`

	// example: 1
	UserID uint64 `json:"userID" sql:"type:integer REFERENCES user(id) ON DELETE CASCADE"`
}

const (
	TaskRunning      = "TaskRunning"
	TaskStarted      = "TaskStarted"
	TaskCreated      = "TaskCreated"
	TaskInitializing = "TaskInitializing"
	TaskStopped      = "TaskStopped"
	TaskError        = "TaskError"
)

func (t *Task) IsRunning() bool {
	return t.Status == TaskStarted || t.Status == TaskInitializing || t.Status == TaskRunning
}

type TaskEnvironment struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:integer REFERENCES task(id) ON DELETE CASCADE"`
	Value  string `json:"value"`
}

type TaskArgument struct {
	ID     uint64 `gorm:"primary_key" json:"id"`
	TaskID uint64 `json:"taskid" sql:"type:integer REFERENCES task(id) ON DELETE CASCADE"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
