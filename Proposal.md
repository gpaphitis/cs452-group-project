## Proposal
### Team Members:
- Pantelis Kanaris
- Giorgos Paphitis

### Problem Definition:
Developing a distributed MapReduce framework in Go, enabling parallel data 
processing across multiple nodes/machines.  
Some of the major challenges and key objectives are:
1. Detection and Handling of errors and crashes across the nodes.
2. Even distribution of work across the nodes.
3. Scaling according to workload.

### Achievable Goals:
Implementing a foundational framework that performs all the previously 
stated objectives, with performance and efficiency considered as secondary
objectives to the implementation.

### Implementation steps:
1. Define the framework’s communication architecture between the 
scheduler and worker nodes.
2. Implementation of the map and reduce operations.
3. Implementation of error detection, handling and recovery.
4. Reasearch and implementation of the workload distribution algorithm.

### Timeline:
- **October:** Complete steps: 1,2
- **November:** Complete steps: 3,4

### Work Assignment:
- **Pantelis Kanaris:** Research and implementation of Map function and the error 
detection.
- **Giorgos Paphitis:** Implementation of Reduce function and error handling
Together: Decisions about implementation and the workload distribution 
algorithm.
