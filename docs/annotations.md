# Annotations

Annotations allow you to modify the default behaviour of the operator.

They can be added to `metadata.annotations` like here:

```
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
  annotations:
    mongodb.com/atlas-resource-policy: keep
spec:
  name: Test Atlas Operator Project
```

### mongodb.com/atlas-resource-policy=keep

If `mongodb.com/atlas-resource-policy` is set to `keep` operator will not delete the Atlas resource when you delete the k8s resource.

### mongodb.com/atlas-reconciliation-policy=skip

If `mongodb.com/atlas-reconciliation-policy` is set to `skip` the operator doesn't start the reconciliation for the resource.

This allows to pause the syncing with the spec for as long as this annotation is added. This might be useful if you want to make manual changes to resource and do not want the operator to undo them. As soon as this annotation is removed the operator should reconcile the resource and sync it back with the spec.