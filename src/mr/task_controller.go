package mr

import "sync"

type TaskState int

const (
	TaskStateProcessing TaskState = iota
	TaskStateCompleted
	TaskStateIdle
)

type TaskType int

const (
	TaskTypeMap TaskType = iota
	TaskTypeReduce
)

type Task struct {
	ID    int
	State TaskState
	Data  []string
}

func NewTask(id int, data []string) *Task {
	return &Task{
		ID:    id,
		State: TaskStateIdle,
		Data:  data,
	}
}

type TaskController struct {
	Tasks []*Task
	mu    sync.Mutex
}

func NewTaskController(input [][]string) *TaskController {
	tc := &TaskController{}

	// get list of tasks
	tasks := make([]*Task, len(input))
	for i, x := range input {
		t := NewTask(i, x)
		tasks[i] = t
	}

	tc.Tasks = tasks

	return tc
}

// check whether tasks are completed
func (tc *TaskController) Done() bool {
	// check that all tasks are completed
	for _, task := range tc.Tasks {
		if task.State == TaskStateIdle || task.State == TaskStateProcessing {
			return false
		}
	}
	return true
}

func (tc *TaskController) getIdleTasks() []*Task {
	var r []*Task
	for _, t := range tc.Tasks {
		if t.State == TaskStateIdle {
			r = append(r, t)
		}
	}
	return r
}

var wg sync.WaitGroup

// finds the next idle task to assign
func (tc *TaskController) NextTask() *Task {
	// lock so we don't assign the same task twice
	tc.mu.Lock()
	defer tc.mu.Unlock()

	for {
		// check if all done or all processing
		if len(tc.getIdleTasks()) == 0 {
			if tc.Done() {
				return nil
			}
			// otherwise wait to finish processing
			wg.Wait()
		}

		idle := tc.getIdleTasks()

		// use first task
		t := idle[0]
		wg.Add(1)
		t.State = TaskStateProcessing
		return t
	}
}
