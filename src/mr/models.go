package mr

// task types enum
type TaskType int

const (
	TaskTypeMap TaskType = iota
	TaskTypeReduce
)

// state of task
type TaskState int

const (
	TaskStateIdle TaskState = iota
	TaskStateInProgress
	TaskStateCompleted
)

// task data structure
type Task struct {
	Input      string
	TaskID     int
	NumOutputs int
	TaskType   TaskType
	Completed  bool
	TaskState  TaskState
}

type KeyValue struct {
	Key   string
	Value string
}
