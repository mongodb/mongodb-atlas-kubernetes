## Running

Example:
```shell
$ go run main.go \
  atlas-api-transformed-v20231115.yaml \
  --category atlas \
  --major-version v20231115 \
  --path '/api/atlas/v2/groups/{groupId}/clusters' --verb post --gvk Cluster.v1.atlas.generated.mongodb.com \
  --path '/api/atlas/v2/groups' --verb post --gvk Group.v1.atlas.generated.mongodb.com \
  --path '/api/atlas/v2/groups/{groupId}/accessList' --verb post --gvk NetworkPermissionEntry.v1.atlas.generated.mongodb.com \
  --path '/api/atlas/v2/groups/{groupId}/customDBRoles/roles' --verb post --gvk UserCustomDBRole.v1.atlas.generated.mongodb.com \
  --output crds-v20231115.yaml
```

```shell
$ go run main.go \
  atlas-api-transformed-v20241113.yaml \
  --category atlas \
  --major-version v20241113 \
  --path '/api/atlas/v2/groups/{groupId}/flexClusters' --verb post --gvk FlexCluster.v1.atlas.generated.mongodb.com \
  --path '/api/atlas/v2/groups/{groupId}/clusters/{clusterName}/fts/indexes' --verb post --gvk SearchIndex.v1.atlas.generated.mongodb.com \
  --output crds-v20241113.yaml
```

```shell
$ kind create cluster
$ kubectl apply --server-side -f crds-v20231115.yaml
$ kubectl apply --server-side -f crds-v20241113.yaml
$ kubectl apply --server-side -f examples.yaml
$ kubectl get atlas -A                        
NAMESPACE   NAME                                            AGE
foo         cluster.atlas.generated.mongodb.com/new-admin   57s

NAMESPACE   NAME                                           AGE
foo         flexcluster.atlas.generated.mongodb.com/flex   57s

NAMESPACE   NAME                                    AGE
foo         group.atlas.generated.mongodb.com/foo   57s

NAMESPACE   NAME                                                     AGE
foo         networkpermissionentry.atlas.generated.mongodb.com/foo   57s

NAMESPACE   NAME                                                        AGE
foo         usercustomdbrole.atlas.generated.mongodb.com/backup-admin   57s
foo         usercustomdbrole.atlas.generated.mongodb.com/new-admin      57s
```

## Acknowledgements

The work is inspired by:
- https://github.com/fybrik/openapi2crd
- https://github.com/ant31/crd-validation
- https://github.com/kubeflow/crd-validation.
