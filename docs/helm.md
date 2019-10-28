# Helm Chart Template Processor

Dash supports generating Resources from Helm Charts.

Helm support functions in two phases:

  1. If the chart name (`.resource_groups[*].resources[*].helm.chart`) is of the form `<string>/<string>`, it is deemed to be from a [Chart Repository](https://helm.sh/docs/chart_repository/). In this case we will first fetch the chart with `helm fetch --untar ...`.
  2. Once the chart is local, we process the chart with `helm template` to generate the resources and cache them on the file system to be reconciled later.

A sample helm chart resource looks like this.

```
- name: Helm Charts
  helm:
    chart: stable/redis
    valueFiles:
    - redis-vars.yaml
    values:
      replicas: 3
      labels:
        something: foo
        else: bar
```
