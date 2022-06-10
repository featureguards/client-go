# Go SDK for FeatureGuards

[![Go Reference](https://pkg.go.dev/badge/github.com/featureguards/featureguards-go/v2)](https://pkg.go.dev/github.com/featueguards/featureguards-go/v2)

The official [FeatureGuards][featureguards] Go client library.

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

```sh
go mod init
```

Then, reference `featureguards-go/v2` in a Go program with `import`:

```go
import (
	featureguards "github.com/featureguards/featureguards-go/v2"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the `featureguards-go` module automatically.

Alternatively, you can also explicitly `go get` the package into a project:

```bash
go get -u github.com/featureguards/featureguards-go/v2
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
```

### IsOn with attributes

```go
// Create the client once.
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// The only require parameter is the API Key. Other options are available. See [goref].
ft := featureguards.New(ctx, featureguards.WithApiKey("API_KEY"), featureguards.WithDefaults(map[string]bool{"TEST": true}))

// Call IsOn multiple times.
on, _ := ft.IsOn("FOO", featureguards.WithAttributes(
	featureguards.Attributes{}.Int64("user_id", 123).String("company_slug", "acme")))
```

[goref]: https://pkg.go.dev/github.com/featureguards/featureguards-go
[issues]: https://github.com/featureguards/featureguards-go/issues/new
[modules]: https://github.com/golang/go/wiki/Modules
[package-management]: https://code.google.com/p/go-wiki/wiki/PackageManagementTools
[pulls]: https://github.com/featureguards/featureguards-go/pulls
[featureguards]: https://featureguards.com
