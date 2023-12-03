# API Reference

## Packages
- [cloud.spaceship.com/v1alpha1](#cloudspaceshipcomv1alpha1)


## cloud.spaceship.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group

### Resource Types
- [ScheduledResource](#scheduledresource)



#### Condition

_Underlying type:_ _string_



_Appears in:_
- [ScheduledStatus](#scheduledstatus)





#### ScheduledResource







| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `cloud.spaceship.com/v1alpha1`
| `kind` _string_ | `ScheduledResource`
| `kind` _string_ | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ScheduledSpec](#scheduledspec)_ |  |
| `status` _[ScheduledStatus](#scheduledstatus)_ |  |


#### ScheduledSpec





_Appears in:_
- [ScheduledResource](#scheduledresource)

| Field | Description |
| --- | --- |
| `in` _string_ |  |
| `content` _string_ |  |


#### ScheduledStatus





_Appears in:_
- [ScheduledResource](#scheduledresource)

| Field | Description |
| --- | --- |
| `condition` _[Condition](#condition)_ |  |


