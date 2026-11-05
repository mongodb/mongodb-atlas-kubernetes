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

If `mongodb.com/atlas-reconciliation-policy` is set to `skip` the operator doesn't start reconciliation for the resource.

This pauses syncing with the spec while the annotation is present. You can use it to make manual changes without AKO reverting them. After removing the annotation, AKO reconciles again and syncs back to spec.

For AtlasProject team handoff: if teams were previously managed in `AtlasProject.spec.teams` and should now be managed outside that list, add `mongodb.com/atlas-reconciliation-policy=skip`, remove `spec.teams`, wait for the skip reconciliation to clear teams from last-applied config, then remove `skip`. When reconciliation resumes with empty `spec.teams` and no last-applied teams, AKO leaves Atlas project team assignments unchanged and does not modify `AtlasTeam.status.projects`. If teams were never managed through `spec.teams`, no migration is needed.
