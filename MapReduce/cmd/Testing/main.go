package main

import (
	"cs452/internal/mapreduce"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	appAddr := "node0:9000" // must match the Application.Listen() address in your main.go

	// Connect to the Application RPC server
	client, err := rpc.Dial("tcp", appAddr)
	if err != nil {
		log.Fatalf("Failed to connect to application RPC at %s: %v", appAddr, err)
	}
	defer client.Close()
	fmt.Printf("Connected to Application RPC at %s\n", appAddr)

	// Prepare the RunTask arguments
	args := &mapreduce.RunTaskArgs{
		JobName:    "inverse-index",
		Files:      []string{"test/data/pg-being_ernest.txt", "test/data/pg-dorian_gray.txt"},
		MapName:    "inverseindex/map",
		ReduceName: "inverseindex/reduce",
	}

	// Call the remote RunTask method
	var reply mapreduce.RunTaskReply
	err = client.Call("Application.RunTask", args, &reply)
	if err != nil {
		log.Fatalf("RunTask RPC call failed: %v", err)
	}

	// Print what we got back
	fmt.Printf("RunTask RPC reply: accepted=%v message=%s\n", reply.Accepted, reply.Message)
}
