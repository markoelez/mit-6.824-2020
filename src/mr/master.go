package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Master struct {
	MapManager *Manager
}

func (m *Master) GiveTask(_ *struct{}, reply *Task) error {
	t := m.MapManager.GetTask()

	if t == nil {
		log.Fatal("DONE WITH PHASE 1")
	}

	reply.Input = t.Input
	reply.TaskID = t.TaskID
	reply.NumOutputs = t.NumOutputs
	reply.Completed = t.Completed
	reply.TaskState = t.TaskState
	reply.TaskType = t.TaskType

	return nil
}

func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (m *Master) Done() bool {
	return false
}

func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	mapManager := NewManager(nReduce, files, TaskTypeMap)
	m.MapManager = mapManager

	m.server()
	return &m
}
