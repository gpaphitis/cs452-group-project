package mapreduce

import (
	"fmt"
	"net"
	"sync"
)

// Master holds all the state that the master needs to keep track of.
// RegisterChannel  notifies the master of workers that have gone idle and are in need of new work.
type Master struct {
	sync.Mutex

	address     string
	doneChannel chan bool
	workers     []string // protected by the mutex

	// Per-task information
	jobName string   // Name of currently executing job
	files   []string // Input files
	nReduce int      // Number of reduce partitions

	shutdown chan struct{}
	l        net.Listener
	stats    []int

	funcSpec FunctionSpec
	mapper_completed []string
	reducer_completed []string
}

func (mr *Master) SetFunctionSpec(spec FunctionSpec) {
	mr.Lock()
	defer mr.Unlock()
	mr.funcSpec = spec
}

func (mr *Master) GetFunctionSpec(_ *struct{}, reply *FunctionSpec) error {
	mr.Lock()
	defer mr.Unlock()
	*reply = mr.funcSpec
	return nil
}

// Register is an RPC method that is called by workers
// workers use the RPC call Register() to register after they have started up to report that they are ready to receive tasks
func (mr *Master) Register(args *RegisterArgs, _ *struct{}) error {
	mr.Lock()
	defer mr.Unlock()
	debug("Register: worker %s\n", args.Worker)
	mr.workers = append(mr.workers, args.Worker)
	mr.nReduce++
	return nil
}

// initializes a new Map/Reduce Master
func newMaster(master string) (mr *Master) {
	mr = new(Master)
	mr.address = master
	mr.shutdown = make(chan struct{})
	mr.nReduce = 0
	mr.doneChannel = make(chan bool)
	return
}

// run executes a mapreduce job on the given number of mappers and reducers.
//
// First, it divides up the input file among the given number of mappers, and
// schedules each task on workers as they become available. Each map task bins
// its output in a number of bins equal to the given number of reduce tasks.
// Once all the mappers have finished, workers are assigned reduce tasks.
//
// When all tasks have been completed, the reducer outputs are merged,
// statistics are collected, and the master is shut down.
//
// Note that this implementation assumes a shared file system.
func (mr *Master) run(jobName string, files []string, schedule func(phase jobPhase), finish func()) {
	mr.jobName = jobName
	mr.files = files
	debug("%s: Starting Map/Reduce task %s\n", mr.address, mr.jobName)

	schedule(mapPhase)
	schedule(reducePhase)
	
	mr.merge()
	
	finish()
	debug("%s: Map/Reduce task completed\n", mr.address)

	mr.doneChannel <- true
}

// Wait blocks until the currently scheduled work has completed.
// This happens when all tasks have scheduled and completed, the final output
// have been computed, and all workers have been shut down.
func (mr *Master) Wait() {
	<-mr.doneChannel
}

// killWorkers cleans up all workers by sending each one a Shutdown RPC.
// It also collects and returns the number of tasks each worker has performed.
func (mr *Master) killWorkers() []int {
	mr.Lock()
	defer mr.Unlock()
	ntasks := make([]int, 0, len(mr.workers))
	for _, w := range mr.workers {
		debug("Master: shutdown worker %s\n", w)
		var reply ShutdownReply
		ok := call(w, "Worker.Shutdown", new(struct{}), &reply)
		if ok == false {
			fmt.Printf("Master: RPC %s shutdown error\n", w)
		} else {
			ntasks = append(ntasks, reply.Ntasks)
		}
	}
	return ntasks
}
