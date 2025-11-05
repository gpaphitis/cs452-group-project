package mapreduce

import (
	"encoding/json"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.

// As with doMap(), it does so either directly or through a worker.
// doReduce() collects corresponding files from each map result , and runs the reduce function on each collection.
// This process produces nReduce result files.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	kvs map[string][]string,
	reduceF func(key string, values []string) string,
) {
	mergeFileName := mergeName(jobName, reduceTaskNumber)
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kvList := make([]KeyValue, 0, len(keys))
	for _, key := range keys {
        kvList = append(kvList, KeyValue{
            Key:   key,
            Value: reduceF(key, kvs[key]),
        })
    }
	
	jsonBytes, err := json.MarshalIndent(kvList, "", "  ")
	if err != nil {
		debug("json marshal failed\n")
	}

	// Write JSON to file
	if err := os.WriteFile(mergeFileName, jsonBytes, 0644); err != nil {
		debug("file open fail:%s\n", mergeFileName)
	}
}
