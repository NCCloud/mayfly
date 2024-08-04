<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://github.com/user-attachments/assets/c7b76ed9-1c02-4172-bd6c-c37cf1cf6269">
  <img alt="logo" width="600"  src="https://github.com/user-attachments/assets/65690736-3dbb-4e92-9d18-83b7a6284bd7">
</picture>

> Mayfly is a Kubernetes operator that enables you to have time-based resources. They creates or deletes on the specified time.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nccloud/mayfly)
![GitHub Release](https://img.shields.io/github/v/release/nccloud/mayfly)
[![Go Reference](https://pkg.go.dev/badge/github.com/NCCloud/mayfly.svg)](https://pkg.go.dev/github.com/NCCloud/mayfly)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/nccloud/mayfly/test.yaml?label=tests)
![GitHub issues](https://img.shields.io/github/issues/nccloud/mayfly)
![GitHub License](https://img.shields.io/github/license/nccloud/mayfly)

## 📖 General Information

### 📄 Summary

The Mayfly Operator allows you to have your resources on your cluster for a temporary time by the given expiration or mayfly create the resources at the time you specified.
It deletes those resources from the cluster, according to the Mayfly expiration annotation that you set to specify how long the resource should remain active. This can be used to create temporary resources, temporary accesses, or simply to keep your cluster organized and tidy. 
Also, It creates the resources you specify at the given time by using `ScheduleResource` custom resource definition. You can also merge these two features together, just to have some resource created in the future and only for a specific amount of time.   

### 🛠 Configuration

Mayfly is an easy-to-use and configurable project that uses resource watches and schedulers to delete your resources at the appropriate time. It is simple to set up and customize.
To specify which resources should be monitored and cleaned up, you can set the `RESOURCES` environment variable to a comma-separated list of `{ApiVersion};{Kind}` as text. This allows you to customize which resources are targeted for cleanup with expiration annotations.

Example:
```bash
export RESOURCES="v1;Secret;ConfigMap,test.com/v1alpha;MyCRD"
```

In the example above, Mayfly will watch `v1` `Secret` and `ConfigMap` resources as well as `test.com/v1alpha` `MyCRD`.

## 🚀 Usage

### Resouce Expiration

Once you have determined which resources you want Mayfly to monitor, you can set the `mayfly.cloud.namecheap.com/expire` annotation on those resources with a duration value or an exact date. In case of duration values, mayfly will calculate the expiration time based on the creation time of the resource. In case of exact date values, mayfly will delete the resource when the exact date has passed. See the examples below for duration and exact date values.

Duration Value Example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  annotations:
    mayfly.cloud.namecheap.com/expire: 30s
spec:
  containers:
    - name: alpine
      image: alpine
      command:
        - sleep
        - infinity
```

Exact Date/Time Example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  annotations:
    mayfly.cloud.namecheap.com/expire: "2024-12-31T00:00:00Z"
spec:
  containers:
    - name: alpine
      image: alpine
      command:
        - sleep
        - infinity
```

> [!TIP]
> mayfly uses [araddon/dateparser](https://github.com/araddon/dateparse) for arbitrary date parsing. See [examples](https://github.com/araddon/dateparse?tab=readme-ov-file#extended-example) for possible date formats.

### Scheduled Resource Creation

The `ScheduledResource` CRD allows you to schedule the creation of an object in the future. This can be combined with the expire annotation, enabling Mayfly to create and remove certain objects for a temporary period in the future.

Example:
```yaml
apiVersion: cloud.namecheap.com/v1alpha1
kind: ScheduledResource
metadata:
  annotations:
    mayfly.cloud.namecheap.com/expire: 60m
  name: example
  namespace: default
spec:
  in: 30m
  content: |
    apiVersion: v1
    kind: Secret
    metadata:
      name: example
      namespace: default
      annotations:
        mayfly.cloud.namecheap.com/expire: 30m
    data:
      .secret-file: dmFsdWUtMg0KDQo=
status:
  condition: Scheduled
```
This feature is particularly useful for setting up temporary resources that are only needed for a short period, reducing clutter and improving the efficiency of resource management.

## 🛳️ Deployment

The easiest and most recommended way to deploy the Mayfly operator to your Kubernetes cluster is by using the Helm chart. To do this, you will need to add our Helm repository and install it from there, providing the RESOURCES environment variable as needed. If you prefer, you can also compile the operator and install it using any method you choose.

Example:
```bash
helm repo add nccloud https://nccloud.github.io/charts
helm install mayfly nccloud/mayfly --set "env[0].name=RESOURCES" --set "env[0].value=v1;Secret" #For only secrets
```

## 🛠 Development

You can easily compile and run the Mayfly operator by following these steps:

1) Create a Kubernetes Cluster or change context for the existing one.

```bash
kind create cluster
```

2) Run the project with the following environment variable.

```bash
export RESOURCES=v1;Secret # Mayfly will begin monitoring secrets in the cluster. For more information, see the configuration section.
go run .
```

## 🏷️ Versioning

We use [SemVer](http://semver.org/) for versioning.
To see the available versions, check the [tags on this repository](https://github.com/nccloud/mayfly/tags).

## ⭐️ Documentation

For more information about the functionality provided by this library, refer to the [GoDoc](http://godoc.org/github.com/nccloud/mayfly) documentation.


## 🤝 Contribution

We welcome contributions, issues, and feature requests!<br />
If you have any issues or suggestions, please feel free to check the [issues page](https://github.com/nccloud/mayfly/issues) or create a new issue if you don't see one that matches your problem. <br>
Also, please refer to our [contribution guidelines](CONTRIBUTING.md) for details.

## 📝 License
All functionalities are in beta and is subject to change. The code is provided as-is with no warranties.<br>
[Apache 2.0 License](./LICENSE)<br>
<br><br>
<img alt="logo" width="75" src="https://avatars.githubusercontent.com/u/7532706" /><br>
Made with <span style="color: #e25555;">&hearts;</span> by [Namecheap Cloud Team](https://github.com/NCCloud)
