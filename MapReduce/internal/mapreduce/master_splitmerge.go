package mapreduce

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

func (mr *Master) FetchReducerResults(jobName string,
	reducers []string,
	kvs *map[string]string) error {

	*kvs = make(map[string]string)
	for i, reducer := range mr.reducer_completed {
		args := new(FetchReduceArgs)
		args.JobName=jobName
		args.ReducerTask=i
		var reduce_results string
		ok := call(reducer, "Worker.FetchReduceResults", args, &reduce_results)
		var items []KeyValue
		json.Unmarshal([]byte(reduce_results), &items)
		for _, kv := range items {
        (*kvs)[kv.Key] = kv.Value
    	}
		if !ok {
			debug("Fetch Reduce Results Failed: RPC %s fetch error\n", reducer)
		}
	}
	return nil
}

//The master calls mr.merge(), which merges all the nReduce files produced by the previous step into a single output.
// merge combines the results of the many reduce jobs into a single output file
// XXX use merge sort
func (mr *Master) merge() {
	debug("Merge phase")
	var kvs map[string]string
	mr.FetchReducerResults(mr.jobName, mr.reducer_completed, &kvs)
	// kvs = make(map[string]string)
	// for i := 0; i < mr.nReduce; i++ {
	// 	p := mergeName(mr.jobName, i)
	// 	debug("Merge: read %s\n", p)
	// 	file, err := os.Open(p)
	// 	if err != nil {
	// 		log.Fatal("Merge: ", err)
	// 	}
	// 	dec := json.NewDecoder(file)
	// 	for {
	// 		var kv KeyValue
	// 		err = dec.Decode(&kv)
	// 		if err != nil {
	// 			break
	// 		}
	// 		kvs[kv.Key] = kv.Value
	// 	}
	// 	file.Close()
	// }
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file, err := os.Create("mrtmp." + mr.jobName)
	if err != nil {
		log.Fatal("Merge: create ", err)
	}
	w := bufio.NewWriter(file)
	for _, k := range keys {
		fmt.Fprintf(w, "%s: %s\n", k, kvs[k])
	}
	w.Flush()
	file.Close()
}

// removeFile is a simple wrapper around os.Remove that logs errors.
func removeFile(n string) {
	err := os.Remove(n)
	if err != nil {
		log.Fatal("CleanupFiles ", err)
	}
}

// CleanupFiles removes all intermediate files produced by running mapreduce.
func (mr *Master) CleanupFiles() {
	for i := range mr.files {
		for j := 0; j < mr.nReduce; j++ {
			removeFile(reduceName(mr.jobName, i, j))
		}
	}
	for i := 0; i < mr.nReduce; i++ {
		removeFile(mergeName(mr.jobName, i))
	}
	removeFile("mrtmp." + mr.jobName)
}
