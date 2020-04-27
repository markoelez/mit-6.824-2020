package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Master struct {
	mc *TaskController
}

func (m *Master) AssignTaskToWorker(args *TaskAssignmentArgs, reply *TaskAssignment) error {
	task := m.mc.NextTask()
	if task != nil {
		reply.ID = task.ID
		reply.Data = task.Data
		reply.Type = TaskTypeMap
		reply.Finished = false
		return nil
	}
	return nil
}

// listen for RPCs from workers
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// check if job is completed
func (m *Master) Done() bool {
	return m.mc.Done()
}

//
// create a Master.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	// format inputs
	map_inputs := make([][]string, len(files))
	for i, f := range files {
		map_inputs[i] = []string{f}
	}

	// get map task controller
	mtc := NewTaskController(map_inputs)

	m := Master{
		mc: mtc,
	}

	m.server()
	return &m
}
