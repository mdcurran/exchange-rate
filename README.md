# Exchange Rate API

Application structure is inspired by
[golang-standards/project-layout](https://github.com/golang-standards/project-layout).

## Development

### Installing Go

For macOS users (such as myself), using Homebrew is recommended:

```
brew install go
```

Otherwise consult the official Go
[installation instructions](https://golang.org/doc/install).

### Build and Run the Go Binary

From the application root directory:

```
$(cd cmd/server && go run main.go)
```

To verify the application is running successfully, execute 
`curl localhost:8080/api/probe`. This should return
`{"message":"Application healthy!"}.`

### Building the Docker Image

To build the `exchange-rate` Docker image; execute from the application root
directory:

```
docker build -t exchange-rate .
```

### Running with Docker

To run the application in a container, whilst mapping port `8080` of the
container to `8080` of the host:

```
docker run -p 8080:8080 exchange-rate
```

Verify the application is running successfully by executing
`curl localhost:8080/api/probe`.

### Running with Kubernetes

Deploying the application on Kubernetes requires having a running cluster with
Helm pre-installed.

For macOS users, [Docker for Desktop](https://docs.docker.com/docker-for-mac/kubernetes/)
is recommended. For everyone else use [microk8s](https://microk8s.io/) or
something similar.

To install [Helm](https://github.com/helm/helm#install) - just follow the
instructions.

Once a Kubernetes cluster is running with Helm configured - run from the
application root:

```
helm install -n exchange-rate k8s/exchange-rate
```

To make the Pod within the cluster accessible:

```
kubectl port-forward svc/exchange-rate 8080:80
```

Verify the application is running successfully by executing
`curl localhost:8080/api/probe`.

Finally, to delete and clean-up the application plus associated Kubernetes
`Deployment` and `Service` objects:

```
helm del --purge exchange-rate
```

## Testing

### Running Unit Tests

To run tests for all packages (ignoring for now there is only one):

```
go test -v ./...
```

_Note: this application uses a multi-stage Docker build to create the image._
_The tests must pass before the image can be built. This is useful to avoid_
_scenarios such as deploying a broken release, but practically much slower_
_for local development. Just run the above command instead._
