## Setup
1. Copy `nfs-setup.sh` to all nodes.
2. For the NFS server run (assume node0)
```bash
sudo ./nfs-setup server
```
3. For the NFS clients run
```bash
sudo ./nfs-setup client node0
```
4. Copy test data files to NFS shared directory
```bash
scp MapReduce/test/data/*.txt gpaphi02@c240g5-110101.wisc.cloudlab.us:/srv/mapReduceData
```
5. Copy docker compose files to all nodes
6. For master run (assume node0)
```bash
docker-compose -f docker/master-compose.yaml up -d
```
7. For workers run
```bash
docker-compose -f docker/worker-compose.yaml up -d
```

## How to Run
## Locally
**Run app:**
```bash
go run cmd/invertedindex/main.go application node0:9000 node0:9100
```

**Run worker:**
```bash
go run cmd/invertedindex/main.go worker node0:9100 node1:7778
```

## Docker
Since our images aren't on DockerHub, we need to copy all code to each node in order to build the image

**Build, run and stop master:**
```bash
docker-compose -f docker/master-compose.yaml build
docker-compose -f docker/master-compose.yaml up -d
```

**Build, run and stop worker:**
```bash
docker-compose -f docker/worker-compose.yaml build
WORKER_PORT=<port> docker-compose -f docker/worker-compose.yaml up -d
```

`WORKER_PORT` if omitted defaults to 7778.  

### Run Multiple Workers on One Machine
To run multiple instances on one machine we have to use a different **port** and **project name** for each.
```bash
# First instance:
docker-compose -f docker/worker-compose.yaml build

# Second instance:
WORKER_PORT=7779 docker-compose -f docker/worker-compose.yaml -p worker_1 up -d
```