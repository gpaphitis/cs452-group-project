package mapreduce

import (
	"fmt"
	"time"
)

// As tasks become available, schedule() decides how to assign those tasks to workers, and how to handle worker failures.
// schedule() starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	var mapper_list []string
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
		mr.mapper_completed = make([]string, ntasks)
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
		mapper_list = mr.mapper_completed
		mr.reducer_completed = make([]string, ntasks)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	stats := make([]bool, ntasks)
	currentWorker := 0

	for {
		count := ntasks
		for i := 0; i < ntasks; i++ {
			if !stats[i] {
				mr.Lock()
				numWorkers := len(mr.workers)
				fmt.Println(numWorkers)
				if numWorkers == 0 {
					mr.Unlock()
					time.Sleep(time.Second)
					continue
				}
				currentWorker = (currentWorker + 1) % numWorkers
				Worker := mr.workers[currentWorker]
				mr.Unlock()
				var file string
				if phase == mapPhase {
					file = mr.files[i]
				}
				args := DoTaskArgs{JobName: mr.jobName, File: file, Phase: phase, TaskNumber: i, NumOtherPhase: nios, FuncSpec: mr.funcSpec, Mappers: mapper_list}
				go func(slot int, worker_ string) {
					success := call(worker_, "Worker.DoTask", &args, new(struct{}))
					if success {
						stats[slot] = true
						if phase == mapPhase {
							mr.mapper_completed[slot] = worker_
							debug("Map task %d completed by %s\n", slot, worker_)
						} else {
							mr.reducer_completed[slot] = worker_
							debug("Reduce task %d completed by %s\n", slot, worker_)
						}
					}
				}(i, Worker)
			} else {
				count--
			}
		}
		if count == 0 {
			break
		}
		time.Sleep(time.Second)
	}
	debug("Schedule: %v phase done\n", phase)
}
