# Go SDK for FeatureGuards

[![Go Reference](https://pkg.go.dev/badge/github.com/featureguards/featureguards-go/v1)](https://pkg.go.dev/github.com/featueguards/featureguards-go/v1)

The official [FeatureGuards][featureguards] Go client library.

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

```sh
go mod init
```

Then, reference featureguards-go/v1 in a Go program with `import`:

```go
import (
	featureguards "github.com/featureguards/featureguards-go/v1"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the featureguards-go module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/featureguards/featureguards-go/v72
```

## Documentation

For details on all the functionality in this library, see the [Go
documentation][goref].

Below are a few simple examples:

### IsOn

```go
// Create the client once.
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// The only require parameter is the API Key. Other options are available. See [goref].
ft := featureguards.New(ctx, featureguards.WithApiKey("API_KEY"), featureguards.WithDefaults(map[string]bool{"TEST": true}))

// Call IsOn multiple times.
on, err := ft.IsOn("TEST")

on, err = ft.IsOn("FOO")
```
