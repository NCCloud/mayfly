# API Reference

## Packages
- [cloud.namecheap.com/v1alpha2](#cloudnamecheapcomv1alpha2)


## cloud.namecheap.com/v1alpha2

Package v1alpha2 contains API Schema definitions for the  v1alpha2 API group

### Resource Types
- [ScheduledResource](#scheduledresource)



#### Condition

_Underlying type:_ _string_





_Appears in:_
- [ScheduledResourceStatus](#scheduledresourcestatus)

| Field | Description |
| --- | --- |
| `Finished` |  |
| `Scheduled` |  |
| `Failed` |  |




#### ScheduledResource









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `cloud.namecheap.com/v1alpha2` | | |
| `kind` _string_ | `ScheduledResource` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ScheduledResourceSpec](#scheduledresourcespec)_ |  |  |  |
| `status` _[ScheduledResourceStatus](#scheduledresourcestatus)_ |  |  |  |


#### ScheduledResourceSpec







_Appears in:_
- [ScheduledResource](#scheduledresource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `schedule` _string_ |  |  |  |
| `content` _string_ |  |  |  |


#### ScheduledResourceStatus







_Appears in:_
- [ScheduledResource](#scheduledresource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `nextRun` _string_ |  |  |  |
| `lastRun` _string_ |  |  |  |
| `condition` _[Condition](#condition)_ |  |  |  |


