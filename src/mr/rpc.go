package mr

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}

func Call(rpcname string, args interface{}, reply interface{}) bool {
	/*
		sockname := masterSock()
		c, err := rpc.DialHTTP("unix", sockname)
	*/
	c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if err != nil {
		log.Fatal("dialing", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
