## My thoughts about Application Local Setup

In my opinion, `docker-compose` or another lightweight container orchestration tool is the best solution for making application up in the local environment.

### Prequisite
Clone the repository by using `git clone git@github.com:HashimovH/softwareengineer-test-task-go.git` command and then `cd softwareengineer-test-task-go`

### Running via Docker
```
docker run -p 8080:8080 ticket-analysis-service-go
```

command will create container in the `8080` port and it will provide gRPC **reflection** since settings contains `APP_ENV=development` inside compose file.


### Using main.go

If you want to start application without containers, after setting up your machine's go settings, you can achieve it by running `go run main.go` in the root folder. Make sure to have `.env` file with the following variable:
```
APP_ENV=development
```

It will provide gRPC reflection and using `gRPCurl` will be easy to test the API.

### K8s
As you can see in the `k8s` folder, we have manifests files which allows to create kubernetes deployment and pod. For local testing, `minikube` is good option to have in order to apply manifests files and interact with the API through Kubernetes.
```
kubectl apply -f k8s/
```