package main

import (
	"cs452/internal/mapreduce"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	appAddr := ":9000" // must match the Application.Listen() address in your main.go

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
		// Files:      []string{"test/data/pg-being_ernest.txt", "test/data/pg-dorian_gray.txt"},
		Files:      []string{"/app/data/pg-being_ernest.txt","/app/data/pg-dracula.txt","/app/data/pg-frankenstein.txt","/app/data/pg-grimm.txt","/app/data/pg-les_miserables.txt","/app/data/pg-moby_dick.txt","/app/data/pg-tale_of_two_cities.txt","/app/data/pg-ulysses.txt","/app/data/pg-dorian_gray.txt","/app/data/pg-emma.txt","/app/data/pg-great_expectations.txt","/app/data/pg-huckleberry_finn.txt","/app/data/pg-metamorphosis.txt","/app/data/pg-sherlock_holmes.txt","/app/data/pg-tom_sawyer.txt","/app/data/pg-war_and_peace.txt"},
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
