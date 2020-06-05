package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Master struct {
	MapManager    *Manager
	ReduceManager *Manager
}

func (m *Master) GiveTask(_ *struct{}, reply *Task) error {
	t := m.MapManager.GetTask()

	if t != nil {
		reply.Input = t.Input
		reply.TaskID = t.TaskID
		reply.NumOutputs = t.NumOutputs
		reply.Completed = t.Completed
		reply.TaskState = t.TaskState
		reply.TaskType = t.TaskType
		return nil
	}

	t = m.ReduceManager.GetTask()

	if t != nil {
		// do reduce tasks
		reply.Input = t.Input
		reply.TaskID = t.TaskID
		reply.NumOutputs = t.NumOutputs
		reply.Completed = t.Completed
		reply.TaskState = t.TaskState
		reply.TaskType = t.TaskType
		return nil
	}

	return nil
}

func (m *Master) SubmitTask(args *Task, _ *struct{}) error {
	if args.TaskType == TaskTypeMap {
		m.MapManager.CompleteTask(args)
	}
	if args.TaskType == TaskTypeReduce {
		m.ReduceManager.CompleteTask(args)
	}
	return nil
}

func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	/*
		sockname := masterSock()
		os.Remove(sockname)
		l, e := net.Listen("unix", sockname)
	*/
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (m *Master) Done() bool {
	return m.MapManager.Done() && m.ReduceManager.Done()
}

func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	mapManager := NewManager(nReduce, files, TaskTypeMap)

	reduceFiles := make([]string, 0)
	for i, _ := range files {
		reduceFiles = append(reduceFiles, fmt.Sprintf("mr-out-%d", i))
	}

	reduceManager := NewManager(nReduce, reduceFiles, TaskTypeReduce)

	m.MapManager = mapManager
	m.ReduceManager = reduceManager

	m.server()
	return &m
}
