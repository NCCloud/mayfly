<img src="https://abload.de/img/mayfly2c7fx3.png" width="600" alt="logo"/>

> Kubernetes operator that allows you to crete ephemeral resources on the cluster that will expire.

## 📖 General Information

### 📄 Summary

Mayfly allows you to have your resources on your cluster for a temporary time.
It deletes those resources from cluster according to mayfly expiration annotation your put to set; how long the resource should be living. When you think about it, you can
use it create temporary resources, temporary accesses or just to keep your cluster always clean and tidy.

### 🛠 Configuration

Mayfly is easy to use and configurable project. Under to hood, it uses resource watches and schedulers to delete your
resource right at the moment when you do.
So, in order to set which resources should be watched and cleanup; you need to set it by the `RESOURCES` environment variable.
This environment variable is comma seperated list of `{ApiVersion};{Kind}` as text.

Example:
```
export RESOURCES="v1;Secret,test.com/v1alpha;MyCRD"
```

## 🚀 Usage
After you successfully set what is being watched by the mayfly, you can set the `mayfly.cloud.spaceship.com/expire` annotation to the resources with a duration value.
Please remember that expiration will be calculated by checking the creation of the resource.

Example:
```
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  annotations:
    mayfly.cloud.spaceship.com/expire: 30s
spec:
  containers:
    - name: alpine
      image: alpine
      command:
        - sleep
        - infinity
```


## 🛳️ Deployment

Easiest, best and recommended way of deploying the mayfly operator into you Kubernetes cluster is using the helm chart of it.
In order make it you should add our helm repository and install it from there by providing the RESOURCES environment variable.
If you don't want to you can also compile easily and install however you want.

Example:
```
helm repo add nccloud UPDATE_THIS
helm install mayfly nccloud/UPDATE_THIS --set RESOURCES="v1;Secret" #For only secrets
```

## 🛠 Development

You can easily compile & run mayfly operator with the following steps.

1) Create a Kubernetes Cluster or change context for the existing one.

```bash
kind create cluster
```

2) Run the project with the following environment variable.

```bash
export RESOURCES=v1;Secret # Mayfly will start to watch secrets in the cluster. Please check configuration section for more.
go run .
```

## 🏷️ Versioning

We use [SemVer](http://semver.org/) for versioning.
For the versions available, see the [tags on this repository](https://github.com/nccloud/mayfly/tags).

## ⭐️ Documentation

For details on all the functionality in this library, see the [GoDoc](http://godoc.org/github.com/nccloud/mayfly) documentation.


## 🤝 Contribution

Contributions, issues and feature requests are welcome!<br />
Feel free to check [issues page](https://github.com/nccloud/mayfly/issues) or create if you find one.


## 📝 License
<img alt="logo" width="75" src="https://avatars.githubusercontent.com/u/7532706" /><br>
This project is [MIT License](https://github.com/nccloud/mayfly) licensed.<br />
Made with <span style="color: #e25555;">&hearts;</span> by DevOps team of Namecheap Cloud.
