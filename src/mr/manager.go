package mr

import (
	"fmt"
	"sync"
)

type Manager struct {
	sync.Mutex

	Cond *sync.Cond

	NumReduce int

	Tasks []*Task
}

func NewManager(numReduce int, inputFiles []string, taskType TaskType) (m *Manager) {
	// split input files into individual tasks
	tasks := make([]*Task, 0)
	if taskType == TaskTypeMap {
		for i, file := range inputFiles {
			t := MakeTask(file, i, numReduce, taskType)
			tasks = append(tasks, t)
		}
	}

	// create manager
	m = new(Manager)
	m.NumReduce = numReduce
	m.Tasks = tasks
	m.Cond = sync.NewCond(m)

	return m
}

func MakeTask(file string, taskID int, numOutputs int, taskType TaskType) *Task {
	return &Task{
		TaskID:     taskID,
		NumOutputs: numOutputs,
		TaskType:   taskType,
		Completed:  false,
		TaskState:  TaskStateIdle,
		Input:      file,
	}
}

func (m *Manager) GetTask() *Task {
	// search through idle tasks and pick one
	m.Lock()
	defer m.Unlock()

	for {
		idleTasks := m.GetIdleTasks()
		fmt.Printf("IDLE: %v\n", idleTasks)
		if len(idleTasks) == 0 {
			//m.Cond.Wait()
			return nil
		}

		// just take first task
		t := idleTasks[0]
		fmt.Println(t)
		// mark as in progress
		t.TaskState = TaskStateInProgress

		//m.Cond.Wait()

		return t
	}
}

func (m *Manager) GetIdleTasks() []*Task {
	r := make([]*Task, 0)
	for _, t := range m.Tasks {
		if t.TaskState == TaskStateIdle {
			r = append(r, t)
		}
	}
	return r
}
