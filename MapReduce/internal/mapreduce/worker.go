package mapreduce

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"fmt"
)

// Worker holds the state for a server waiting for DoTask or Shutdown RPCs
type Worker struct {
	sync.Mutex
	registry *Registry
	name   string
	Map    func(string, string) []KeyValue
	Reduce func(string, []string) string
	nRPC   int // protected by mutex
	nTasks int // protected by mutex
	l      net.Listener
	master string
}

// DoTask is called by the master when a new task is being scheduled on this
// worker.
func (wk *Worker) DoTask(arg *DoTaskArgs, _ *struct{}) error {
	debug("%s: given %v task #%d on file %s (nios: %d)\n",
		wk.name, arg.Phase, arg.TaskNumber, arg.File, arg.NumOtherPhase)

	switch arg.Phase {
	case mapPhase:
		mapFn, ok := wk.registry.GetMap(arg.FuncSpec.MapName)
        if !ok || mapFn == nil {
            return fmt.Errorf("worker %s: map function %q not found or nil", wk.name, arg.FuncSpec.MapName)
		}
		doMap(arg.JobName, arg.TaskNumber, arg.File, arg.NumOtherPhase, mapFn)
	case reducePhase:
		redFn, ok := wk.registry.GetReduce(arg.FuncSpec.ReduceName)
        if !ok || redFn == nil {
            return fmt.Errorf("worker %s: reduce function %q not found or nil", wk.name, arg.FuncSpec.ReduceName)
		}
		var kvs map[string][]string
		err:=wk.FetchResultsFromMappers(arg.JobName, arg.TaskNumber, arg.Mappers, &kvs)
		if err != nil {
			debug("Fetching Interm results failed")
			return nil
		}
		doReduce(arg.JobName, arg.TaskNumber, kvs, redFn)
	}

    debug("%s: %v task #%d done\n", wk.name, arg.Phase, arg.TaskNumber)
    return nil
}

func (wk *Worker) FetchResultsFromMappers(jobName string,
	reduceTaskNumber int,
	mappers []string,
	kvs *map[string][]string) error {

	*kvs = make(map[string][]string)
	for i, mapper := range mappers {
		args := new(FetchIntermediateArgs)
		args.JobName=jobName
		args.ReducerId=reduceTaskNumber
		args.MapTaskNumber=i
		var intermediate_results string
		ok := call(mapper, "Worker.FetchIntermediateResults", args, &intermediate_results)
		var items []KeyValue
		json.Unmarshal([]byte(intermediate_results), &items)
		for _, item := range items {
			k := item.Key
			v := item.Value
			(*kvs)[k] = append((*kvs)[k], v)
		}
		if !ok {
			debug("Fetch Results Failed: RPC %s fetch error\n", mapper)
		}
	}
	return nil
}

func (wk *Worker) FetchIntermediateResults(arg *FetchIntermediateArgs, reply *string) error {
	filename := reduceName(arg.JobName, arg.MapTaskNumber, arg.ReducerId)
	data, err := os.ReadFile(filename)
	if err != nil {
		debug("Error reading file %s\n", filename)
		return err
	}
	*reply = string(data)
	return nil
}

func (wk *Worker) FetchReduceResults(arg *FetchReduceArgs, reply *string) error {
	filename := mergeName(arg.JobName, arg.ReducerTask)
	data, err := os.ReadFile(filename)
	if err != nil {
		debug("Error reading file %s\n", filename)
		return err
	}
	*reply = string(data)
	return nil
}

// Shutdown is called by the master when all work has been completed.
// We should respond with the number of tasks we have processed.
func (wk *Worker) Shutdown(_ *struct{}, res *ShutdownReply) error {
	debug("Shutdown %s\n", wk.name)
	wk.Lock()
	defer wk.Unlock()
	res.Ntasks = wk.nTasks
	wk.nRPC = 1
	wk.nTasks-- // Don't count the shutdown RPC
	return nil
}

// Tell the master we exist and ready to work
func (wk *Worker) register(master string) {
	wk.master = master
	args := new(RegisterArgs)
	args.Worker = wk.name
	ok := call(master, "Master.Register", args, new(struct{}))
	if ok == false {
		debug("Register: RPC %s register error\n", master)
	}
}

// RunWorker sets up a connection with the master, registers its address, and
// waits for tasks to be scheduled.
func RunWorker(MasterAddress string, me string, nRPC int) {
	debug("RunWorker %s\n", me)
	wk := new(Worker)
	wk.name = me
	//TODO: MAKE THE WORKER POPULATE THE REGISTRY
	wk.registry = NewRegistry()
	wk.registry.PopulateRegistry()
	wk.nRPC = nRPC
	rpcs := rpc.NewServer()
	rpcs.Register(wk)
	os.Remove(me) // only needed for "tcp"
	l, e := net.Listen("tcp", me)
	if e != nil {
		log.Fatal("RunWorker: worker ", me, " error: ", e)
	}
	wk.l = l
	wk.register(MasterAddress)

	// DON'T MODIFY CODE BELOW
	for {
		wk.Lock()
		if wk.nRPC == 0 {
			wk.Unlock()
			break
		}
		wk.Unlock()
		conn, err := wk.l.Accept()
		if err == nil {
			wk.Lock()
			wk.nRPC--
			wk.Unlock()
			go rpcs.ServeConn(conn)
			wk.Lock()
			wk.nTasks++
			wk.Unlock()
		} else {
			break
		}
	}
	wk.l.Close()
	debug("RunWorker %s exit\n", me)
}
