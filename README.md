contextops
----------
[godoc](https://pkg.go.dev/github.com/DanLavine/contexops)


ContextOPS is a helper package to manage common context operation not provided by the
standard Golang [context](https://pkg.go.dev/context) package.

# MergeForDone

MergeForDone can be used to merge any number of contexts into one when you only care
about the `Done()` channel being closed on any number of contexts. This can be useful
in cases where a long or blocking API is taking place. But a server shutdown or client
disconnect can also cancel the API call.
