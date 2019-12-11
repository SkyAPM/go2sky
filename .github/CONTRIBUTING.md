# How to contribute

Contributing to Go2Sky is as easy as opening a pull request with the changes. It is often a good idea to initiate a discussion first by
opening a GitHub issue, to get early feedback and discuss the design of new features.

Every contribution is expected to have proper unit tests and to follow the following code conventions:

### Code conventions

All Golang code checked into this repo must conform to:
- conventions outlined in [Effective Go](https://golang.org/doc/effective_go.html)
- [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) output
  > Note that `goimports` is a strict superset of `go fmt`.

  Import blocks must have three sections: golang standard library imports, third party imports, and first party imports (other Go2Sky code). Each section must be separated from the others by a blank line.
  This can be enforced by running it as follows (make sure you configure your IDE to automatically use this): `goimports -local github.com/SkyAPM/go2sky`
  
- [golangci-lint](https://github.com/golangci/golangci-lint). All linter warnings are treated as errors and code with linter warnings may not be checked in. Supressions of linter warnings is allowed at the discretion of a PR's reviewers.

### Code Reviews

The Go2Sky community will review your pull request before it is merged. This process can take a while, so please be patient and
make sure the pull request is _small and focused_ so reviews can be provided in a timely manner. 

During the review process you may be asked to make some changes to your submission. While working through feedback, it can be beneficial
to create new commits so the incremental change is obvious. This can also lead to a complex set of commits, and having an atomic change
per commit is preferred in the end. Use your best judgement and work with your reviewer as to when you should revise a commit or
create a new one.  

A pull request is considered ready to be merged once it gets at lease one +1 from a reviewer. Once all the changes have been completed
and the pull request is accepted, it must be rebased to the latest upstream version and all tests and lint checks must pass.
