# MongoDB Atlas Operator
Welcome to MongoDB Atlas Operator - the Operator allowing to manage Atlas Clusters from Kubernetes.
> Current Status: *pre-alpha*. We are currently working on the initial set of features that will give users the opportunity to provision Atlas projects, clusters and database users using Kubernetes Specifications and bind connection information into the applications deployed to Kubernetes.   
## Quick Start guide
### Step 1. Deploy Kubernetes operator using all in one config file
```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/deploy/all-in-one.yaml
```
### Step 2. Create Atlas Cluster

**1.** Create Atlas API Key Secret
In order to work with Atlas Operator you'll need to provide the [authentication information](https://docs.atlas.mongodb.com/configure-api-access) that would allow the Atlas Operator communicate with Atlas API. You need to create the secret first:
```
kubectl create secret generic my-atlas-key \
         --from-literal="orgId=<the_atlas_organization_id>" \
         --from-literal="publicApiKey=<the_atlas_api_public_key>" \
         --from-literal="privateApiKey=<the_atlas_api_private_key>"    
```
**2.** Create `AtlasProject` Custom resource

`AtlasProject` represents Atlas Project in our Kubernetes cluster. 
Note: Property `connectionSecretRef` should reference the Secret `my-atlas-key` that was created in previous step:
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
EOF
```
**3.** Create the `AtlasCluster` Resource.
Example below is a minimal configuration for Atlas MongoDB Cluster in AWS US East region, size M10. For a full list of properties check 
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
You can use the following command to track the status of the Atlas Cluster Resource
```
kubectl get atlasclusters.atlas.mongodb.com -o yaml
```

