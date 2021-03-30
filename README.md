# MongoDB Atlas Operator
[![MongoDB Atlas Operator](https://github.com/mongodb/mongodb-atlas-kubernetes/workflows/Test/badge.svg)](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/test.yml?query=branch%3Amain)
[![MongoDB Atlas Go Client](https://img.shields.io/badge/Powered%20by%20-go--client--mongodb--atlas-%2313AA52)](https://github.com/mongodb/go-client-mongodb-atlas)

Welcome to the MongoDB Atlas Operator - a Kubernetes Operator which manages MongoDB Atlas Clusters from Kubernetes.

> Current Status: *beta*. The Operator gives users the ability to provision
> Atlas projects, clusters and database users using Kubernetes Specifications and bind connection information
> into applications deployed to Kubernetes. More features like private endpoints, backup management, LDAP/X.509 authentication, etc. 
> are not supported yet.

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
kubectl create secret generic mongodb-atlas-operator-api-key \
         --from-literal="orgId=<the_atlas_organization_id>" \
         --from-literal="publicApiKey=<the_atlas_api_public_key>" \
         --from-literal="privateApiKey=<the_atlas_api_private_key>" \
         -n mongodb-atlas-system
```

**2.** Create an `AtlasProject` Custom Resource

The `AtlasProject` CustomResource represents Atlas Projects in our Kubernetes cluster. You need to specify 
`projectIpAccessList` with the IP addresses or CIDR blocks of any hosts that will connect to the Atlas Cluster.
```
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: Test Atlas Operator Project
  projectIpAccessList:
    - ipAddress: "192.0.2.15"
      comment: "IP address for Application Server A"
    - ipAddress: "203.0.113.0/24"
      comment: "CIDR block for Application Server B - D"
EOF
```
**3.** Create an `AtlasCluster` Custom Resource.
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

**4.** Create a database user password Kubernetes Secret
```
kubectl create secret generic the-user-password --from-literal="password=P@@sword%"
```

**5.** Create an `AtlasDatabaseUser` Custom Resource
In order to connect to an Atlas Cluster the database user needs to be created. `AtlasDatabaseUser` resource should reference
the Kubernetes Secret created in the previous step.
```
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasDatabaseUser
metadata:
  name: my-database-user
spec:
  roles:
    - roleName: "readWriteAnyDatabase"
      databaseName: "admin"
  projectRef:
    name: my-project
  username: theuser
  passwordSecretRef:
    name: the-user-password
EOF
```
**6.** Wait for the `AtlasDatabaseUser` Custom Resource to be ready

Wait until the AtlasDatabaseUser resource gets to "ready" status (it will wait until the cluster is created that may take around 10 minutes):
```
kubectl get atlasdatabaseusers my-database-user -o=jsonpath='{.status.conditions[?(@.type=="Ready")].status}'
True
```
### Step 4. Connect your application to the Atlas Cluster

The Atlas Operator will create a Kubernetes Secret with the information necessary to connect to the Atlas Cluster created
in the previous step. It can be used by your application in the same Kubernetes Cluster to connect to 
**1.** Create an `AtlasDatabaseUser` Resource

### Step 3. Inspect Atlas Cluster Status
You can use the following command to check the status of the Atlas Cluster Resource
```
kubectl get atlasclusters -o yaml
```
