## Setup
1. Copy `nfs-setup.sh.sh` to all nodes.
```bash
parallel-scp -h hosts.txt -l gpaphi02 nfs-setup.sh /users/gpaphi02/
```
2. For the NFS server run (assume node0)
```bash
sudo ./nfs-setup.sh server
```
3. For the NFS clients run
```bash
sudo ./nfs-setup.sh client node0
```
4. Copy test data files to NFS shared directory
```bash
scp MapReduce/test/data/*.txt gpaphi02@c240g5-110101.wisc.cloudlab.us:/srv/mapReduceData
```
5. Copy docker compose files to all nodes
```bash
parallel-scp -r -h hosts.txt -l gpaphi02 docker /users/gpaphi02/
```
6. For master run (assume node0)
```bash
docker-compose -f docker/master-compose.yaml up -d
```
7. For workers run
```bash
docker-compose -f docker/worker-compose.yaml up -d
```

## How to Run
## Docker
### Build
```bash
# Master
docker build --tag giorgospaphitis/personal:cs452-worker -f docker/worker.Dockerfile .
docker push giorgospaphitis/personal:cs452-master

# Worker
docker build --tag giorgospaphitis/personal:cs452-worker -f docker/worker.Dockerfile .
docker push giorgospaphitis/personal:cs452-worker
```

**Run and stop master:**
```bash
# Master
docker-compose -f docker/master-compose.yaml up -d
docker-compose -f docker/master-compose.yaml down

# Worker
WORKER_PORT=<port> docker-compose -f docker/worker-compose.yaml up -d
docker-compose -f docker/worker-compose.yaml down
```
`WORKER_PORT` if omitted defaults to 7778.  

### Run Multiple Workers on One Machine
To run multiple instances on one machine we have to use a different **port** and **project name** for each.
```bash
# First instance:
docker-compose -f docker/worker-compose.yaml up -d
docker-compose -f docker/worker-compose.yaml down

# Second instance:
WORKER_PORT=7779 docker-compose -f docker/worker-compose.yaml -p worker_1 up -d
docker-compose -f docker/worker-compose.yaml -p worker_1 down
```

## Locally
**Run app:**
```bash
go run cmd/invertedindex/main.go application node0:9000 node0:9100
```

**Run worker:**
```bash
go run cmd/invertedindex/main.go worker node0:9100 node1:7778
```



## Test Inverse Index Result
Compare the following ouputs to he ones from the lab
```bash
head -n5 /mnt/mapReduceData/mrtmp.inverse-index
A: 16 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
ABC: 2 pg-les_miserables.txt,pg-war_and_peace.txt
ABOUT: 2 pg-moby_dick.txt,pg-tom_sawyer.txt
ABRAHAM: 1 pg-dracula.txt
ABSOLUTE: 1 pg-les_miserables.txt

sort -k1,1 /mnt/mapReduceData/mrtmp.inverse-index. | sort -snk2,2 /mnt/mapReduceData/mrtmp.inverse-index | grep -v '16' | tail -10
women: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
won: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
wonderful: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
words: 15 pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
worked: 15 pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
worse: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
wounded: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
yes: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
younger: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
yours: 15 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
```