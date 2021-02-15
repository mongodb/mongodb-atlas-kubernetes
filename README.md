# MongoDB Atlas Operator
![MongoDB Atlas Operator](https://github.com/mongodb/mongodb-atlas-kubernetes/workflows/Test/badge.svg)
[![MongoDB Atlas Go Client](https://img.shields.io/badge/Powered%20by%20-go--client--mongodb--atlas-%2313AA52)](https://github.com/mongodb/go-client-mongodb-atlas)

Welcome to the MongoDB Atlas Operator - a Kubernetes Operator which manages MongoDB Atlas Clusters from Kubernetes.

> Current Status: *alpha*. We are currently working on the initial set of features that will give users the opportunity
> to provision Atlas projects, clusters and database users using Kubernetes Specifications and bind connection information
> into the applications deployed to Kubernetes.

## Quick Start guide
### Step 1. Deploy Kubernetes operator using all in one config file
```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/main/deploy/all-in-one.yaml
```
### Step 2. Create Atlas Cluster

**1.** Create an Atlas API Key Secret
In order to work with the Atlas Operator you need to provide [authentication information](https://docs.atlas.mongodb.com/configure-api-access)
 to allow the Atlas Operator to communicate with Atlas API. Once you have generated a Public and Private key in Atlas, you can create a Kuberentes Secret with:
```
kubectl create secret generic my-atlas-key \
         --from-literal="orgId=<the_atlas_organization_id>" \
         --from-literal="publicApiKey=<the_atlas_api_public_key>" \
         --from-literal="privateApiKey=<the_atlas_api_private_key>"
```
**2.** Create the `AtlasProject` Custom Resource

The `AtlasProject` CustomResource represents Atlas Projects in our Kubernetes cluster.
Note: the property `connectionSecretRef` should reference the Secret which holds your Atlas API keys that was created in the previous step (`my-atlas-key` in the example above).
```
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: Test Atlas Operator Project
  connectionSecretRef:
    name: my-atlas-key
  projectIpAccessList:
    - ipAddress: "192.0.2.15"
      comment: "IP address for Application Server A"
    - ipAddress: "203.0.113.0/24"
      comment: "CIDR block for Application Server B - D"
EOF
```
**3.** Create an `AtlasCluster` Resource.
The example below is a minimal configuration to create an M10 Atlas cluster in the AWS US East region. For a full list of properties, check
`atlasclusters.atlas.mongodb.com` [CRD specification](config/crd/bases/atlas.mongodb.com_atlasclusters.yaml)):
```
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasCluster
metadata:
  name: my-atlas-cluster
spec:
  name: "Test-cluster"
  projectRef:
    name: my-project
  providerSettings:
    instanceSizeName: M10
    providerName: AWS
    regionName: US_EAST_1
EOF
```
### Step 3. Inspect Atlas Cluster Status
You can use the following command to check the status of the Atlas Cluster Resource
```
kubectl get atlasclusters.atlas.mongodb.com -o yaml
```
