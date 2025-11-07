# How to Run
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
HOST_PORT=<port> docker-compose -f docker/worker-compose.yaml up -d
```

`HOST_PORT` if omitted defaults to 7778.  

### Run Multiple Workers on One Machine
To run multiple instances on one machine we have to use a different **port** and **project name** for each.
```bash
# First instance:
docker-compose -f docker/worker-compose.yaml build

#Second instance:
HOST_PORT=7779 docker-compose -f docker/worker-compose.yaml -p worker_1 up -d
```