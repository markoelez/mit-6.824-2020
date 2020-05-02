package mr

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

type ByKey []KeyValue

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func Worker(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	intermediate := []KeyValue{}

	for {

		t := getTask()
		fmt.Printf("CURRENT TASK: %v\n", t)

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
			intermediate = append(intermediate, kva...)

			submitTask(t)
		} else if t.TaskType == TaskTypeReduce {
			sort.Sort(ByKey(intermediate))

			oname := fmt.Sprintf("mr-out-%d", t.TaskID)
			ofile, _ := os.Create(oname)
			w := bufio.NewWriter(ofile)

			i := 0
			for i < len(intermediate) {
				j := i + 1
				for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
					j++
				}

				values := []string{}
				for k := i; k < j; k++ {
					values = append(values, intermediate[k].Value)
				}

				output := reducef(intermediate[i].Key, values)
				fmt.Fprintf(w, "%v %v\n", intermediate[i].Key, output)

				i = j
			}
			ofile.Close()
			submitTask(t)
		}

		time.Sleep(time.Second * 1)
	}
}

func submitTask(t *Task) {
	Call("Master.SubmitTask", &t, &struct{}{})
}

func getTask() *Task {
	t := Task{}
	Call("Master.GiveTask", &struct{}{}, &t)
	return &t
}
