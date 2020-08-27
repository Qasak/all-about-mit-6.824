package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "fmt"


type Master struct {
	// Your definitions here.
	files map[string]bool


}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (m *Master) SelectOneFile() string {
	for k, v := range m.files {
		if v == false {
			m.files[k]=true
			return k
		}
	}
	return ""
}
func (m *Master) SendTask(args *WorkerArgs, reply *WorkerReply) error {
	if args.Request == "get_task" {
		reply.Filename = m.SelectOneFile()
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	// -----------test point-----------
	fmt.Println("sockname:", sockname)
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	// // -----------test point-----------
	// fmt.Println("l:", l)
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	// -----------test point-----------
	// fmt.Println(files)
	m := Master{}
	m.files=make(map[string]bool)
	for _, file:=range files {
		m.files[file]=false
	}
	// Your code here.


	m.server()
	return &m
}
