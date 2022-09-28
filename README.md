<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://abload.de/img/mayfly-logo-lightm5ib0.png">
  <img alt="logo" width="600"  src="https://abload.de/img/mayfly-logo-darkt9eye.png">
</picture>

> Mayfly is a Kubernetes operator that enables you to create temporary resources on the cluster that will expire after a certain period of time.

## 📖 General Information

### 📄 Summary

The Mayfly Operator allows you to have your resources on your cluster for a temporary time.
It deletes those resources from the cluster, according to the Mayfly expiration annotation that you set to specify how long the resource should remain active. This can be used to create temporary resources, temporary accesses, or simply to keep your cluster organized and tidy.

### 🛠 Configuration

Mayfly is an easy-to-use and configurable project that uses resource watches and schedulers to delete your resources at the appropriate time. It is simple to set up and customize.
To specify which resources should be monitored and cleaned up, you can set the `RESOURCES` environment variable to a comma-separated list of `{ApiVersion};{Kind}` as text. This allows you to customize which resources are targeted for cleanup.

Example:
```
export RESOURCES="v1;Secret,test.com/v1alpha;MyCRD"
```

## 🚀 Usage
Once you have determined which resources you want Mayfly to monitor, you can set the `mayfly.cloud.namecheap.com/expire` annotation on those resources with a duration value. This will cause Mayfly to delete the resources once the specified duration has passed, based on the time of their creation.
Keep in mind that the expiration will be calculated based on the creation time of the resource.

Example:
```
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


## 🛳️ Deployment

The easiest and most recommended way to deploy the Mayfly operator to your Kubernetes cluster is by using the Helm chart. To do this, you will need to add our Helm repository and install it from there, providing the RESOURCES environment variable as needed. If you prefer, you can also compile the operator and install it using any method you choose.

Example:
```
helm repo add nccloud https://nccloud.github.io/charts
helm install mayfly nccloud/mayfly --set RESOURCES="v1;Secret" #For only secrets
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
Made with <span style="color: #e25555;">&hearts;</span> by the DevOps team of [Namecheap Cloud](https://github.com/NCCloud)
