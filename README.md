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

## Testing

### Running Unit Tests

To run tests for all packages (ignoring for now there is only one):

```
go test -v ./...
```
