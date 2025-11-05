[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/ppZxTnKS)



run worker:
 go run cmd/invertedindex/main.go worker node0:9100 node1:7778

 run app:
 go run cmd/invertedindex/main.go application node0:9000 node0:9100
