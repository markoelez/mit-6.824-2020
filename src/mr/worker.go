package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
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

		// process task
		if t.TaskType == TaskTypeMap {
			file, err := os.Open(t.Input)
			if err != nil {
				log.Fatalf("cannot open %v\n", t.Input)
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalf("cannot open %v\n", t.Input)
			}
			file.Close()
			kva := mapf(t.Input, string(content))

			oname := fmt.Sprintf("mr-out-%d", t.TaskID)
			fmt.Printf("encoding output file %s for task file %s\n", oname, t.Input)
			ofile, _ := os.Create(oname)

			enc := json.NewEncoder(ofile)
			for _, kv := range kva {
				err := enc.Encode(&kv)
				if err != nil {
					log.Fatal("Error encoding output file")
				}
			}

		}

		time.Sleep(time.Second * 2)
	}
}

func getTask() *Task {
	t := Task{}
	Call("Master.GiveTask", &struct{}{}, &t)
	return &t
}
