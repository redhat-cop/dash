# Dash Prototype

Dash is an automation framework for Kubernetes. The aim is to combine the power of [Kubernetes declarative management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/) and popular templating engines to provide a full-featured automation strategy for Kubernetes cluster and application management.

Dash helps to enable GitOps and Infrastructure as Code workflows in enterprises with complex use cases.

## Features

### Template Support

Currently Dash supports the following templates:

* FileTemplates - not really templates, just static resource definitions
* [HelmCharts](./docs/helm.md)

## Quickstart

To run, first log into your kubernetes cluster, and:

```
export GO111MODULE=on
go mod vendor
go run cmd/dash.go -i examples/default/
```

## Running Tests

There is some automated test coverage in the libraries. You can run all tests with:

```
go test ./...
```
