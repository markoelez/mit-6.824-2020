package mr

import (
	"fmt"
	"hash/fnv"
	"time"
)

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	for {

		t := getTask()

		fmt.Printf("TASK: %v\n", t)

		time.Sleep(time.Second * 2)
	}
}

func getTask() *Task {
	t := Task{}
	Call("Master.GiveTask", &struct{}{}, &t)
	return &t
}

func CallExample() {

	args := ExampleArgs{
		X: 99,
	}

	reply := ExampleReply{}

	Call("Master.Example", &args, &reply)

	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}
