package mr

type TaskAssignmentArgs struct{}
type TaskOutputArgs struct{}

// get task from master server
type TaskAssignment struct {
	ID       int
	Data     []string
	Type     TaskType
	Finished bool
}

type TaskOutput struct {
	ID   int
	Data []string
	Type TaskType
}
