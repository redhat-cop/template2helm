# Development and Contribution to template2helm

We encourage contributions to this project! We follow a fairly standard forking workflow for open source projects. This document provides some information about getting your environment set up.

In general, the requirements to contribute to this project are as follows:

- A [Go 1.12](https://golang.org/dl/) enviornment
- [yq](https://pypi.org/project/yq/)
- (optional) A kubernetes cluster like [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

## Build Binary

The default target in the makefile will build the project binary in the local directory. From there you can manually test it.

```
$ make

$ ./template2helm
Template2helm converts an OpenShift Template into a Helm Chart.
      For more info, check out https://github.com/redhat-cop/template2helm

Usage:
  template2helm [command]

Available Commands:
  convert     Given the path to an OpenShift template file, spit out a Helm chart.
  help        Help about any command
  version     Print the version number of template2helm

Flags:
  -h, --help   help for template2helm

Use "template2helm [command] --help" for more information about a command.
```

## Running Tests

There is some automated test coverage in the libraries. You can run all tests with:

```
$ make test_e2e
```

## Cutting Releases

We use the [github-release](https://github.com/aktau/github-release) utility to automate creating releases of our project. In order to cut a release, a repository admin must first generate a [GitHub API token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line). Then run the following:

```
# Create a new tag for the release
git tag -a <version> -m "Release <version>"
export GITHUB_TOKEN=...
go get github.com/aktau/github-release
make release
```
