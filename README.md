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
