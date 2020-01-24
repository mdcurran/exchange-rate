# Exchange Rate API

Application structure is inspired by
[golang-standards/project-layout](https://github.com/golang-standards/project-layout).

To get the exchange rate of a single currency, and a naïve recommendation of
whether to buy the currency or not:

```
curl -X GET 'localhost:8080/api/rate?currency=USD'
```

After getting the application running using one of the steps below. The
currency can also be changed to be USD, GBP + this will work for any other
currency supported on exchangeratesapi.io.

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

## General Thought and Potential Improvements

Some commentary about a couple of design decisions, and ideas about what I'd
improve if there was more time.

### Test More

Ideally I'd have written tests for `handleRate()`, `buildURI()` and so on. But
I think being able to properly mock out the call to the third-party API,
restructuring the application appropriately, and thinking of all the possible
scenarios would've taken far far too much time.

### Kubernetes

The Helm chart I wrote is really bare-bones. I didn't include things like:

* Configuration through a `values.yaml` file.
* Multiple replicas of the application.
* Better labelling, and so on.

Again this is more a time thing.

### Custom Errors

Deployable applications should have some kind of structured error output. Even
in this application, if it's just an error string + an appropriate HTTP status
code encoded in JSON.

I think this is the bare minimum I'd feel confident with when deploying an
application to production.

### scratch Docker Image

I've used a multi-stage Docker build, with the final image being based on the
minimal image `scratch`. This has a few benefits:

1. Smaller final image size, maximising the number of replicas that could run
   on a host.
2. Faster image pulls from a repository. Will speed up how quickly the
   application can be deployed.
3. Limited attack vectors for malicious actors.

While these are really good reasons for using `scratch`, I've come across
situations where it might not be the best solution. For example, vulnerability
scanning on AWS ECR doesn't work on `scratch`-based images. Depending on an
organisation's security best practices this might be OK, it might not be.
