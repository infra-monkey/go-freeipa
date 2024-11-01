# go-freeipa

[![Go Reference](https://pkg.go.dev/badge/github.com/infra-monkey/go-freeipa.svg)](https://pkg.go.dev/github.com/infra-monkey/go-freeipa)

A generated golang client for the FreeIPA API.

See the [godoc](https://godoc.org/github.com/infra-monkey/go-freeipa/freeipa) for
usage examples and more info. See
[`developing.md`](https://github.com/infra-monkey/go-freeipa/blob/master/developing.md)
for information on how this library is generated.

## Features

* Kerberos and password API authentication

## History and credits

### Origins

This project was initially created by [Philippe Voinov](https://github.com/tehwalris).

The initial project is now _archived_ and [unmaintained](https://github.com/tehwalris/go-freeipa#unmaintained).

We want to thank [Philippe Voinov](https://github.com/tehwalris) for all the awesome work he did on this library.

### Future

[IN2P3 Computing Centre](https://cc.in2p3.fr) will provide support for this library and try to make it evolve as much as possible.

Feel free to open [issues](https://github.com/infra-monkey/go-freeipa/issues) or [pull requests](https://github.com/infra-monkey/go-freeipa/pulls) to discuss.

## Known issues

### pkg.go.dev website refuses to display documentation

[![pkg.go.dev license restrictions](doc/go_freeipa_doc_license_issue.png)](https://pkg.go.dev/github.com/infra-monkey/go-freeipa/freeipa)

The official [pkg.go.dev license policy](https://pkg.go.dev/license-policy) does not actually accepts the [CeCILL LICENSE](https://cecill.info/index.en.html) we're using.

We're aware of this and an [issue was opened to request the recognition of the CeCILL LICENSE on the pkg.go.dev](https://github.com/golang/go/issues/52126).
