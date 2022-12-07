# monitoring-operator

Reconciles Kubernetes Resources against different monitoring tool REST APIs.

## Description

There are monitoring tools without integration with the Kubernetes API. Due to that, since they only offer a REST API,
this operator serves as a shim where we can create Kubernetes Resources which will be reconciled against these REST
APIs, removing the hassle from the infrastructure admin.
