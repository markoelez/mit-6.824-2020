package mr

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"time"
)

type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// entrypoint
func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	for {
		ta := getTask()
		fmt.Printf("TASK %v\n", ta)
		time.Sleep(4 * time.Second)
	}
}

func getTask() TaskAssignment {
	args := TaskAssignmentArgs{}
	reply := TaskAssignment{}

	call("Master.AssignTaskToWorker", &args, &reply)

	return reply
}

// sends RPC
func call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
