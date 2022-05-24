# MongoDB Atlas Operator (Beta)

[![MongoDB Atlas Operator](https://github.com/mongodb/mongodb-atlas-kubernetes/workflows/Test/badge.svg)](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/test.yml?query=branch%3Amain)
[![MongoDB Atlas Go Client](https://img.shields.io/badge/Powered%20by%20-go--client--mongodb--atlas-%2313AA52)](https://github.com/mongodb/go-client-mongodb-atlas)

The MongoDB Atlas Operator provides a native integration between the Kubernetes orchestration platform and MongoDB Atlas
— the only multi-cloud document database service that gives you the versatility you need to build sophisticated and
resilient applications that can adapt to changing customer demands and market trends.

> Current Status: *Beta*.
> The Operator gives users the ability to provision
> Atlas projects, clusters and database users using Kubernetes Specifications and bind connection information
> into applications deployed to Kubernetes or via Private Endpoints on AWS or Azure.
> More features like private endpoints, backup management, LDAP/X.509 authentication, etc. are yet to come.

The full documentation for the Operator can be found [here](https://docs.atlas.mongodb.com/atlas-operator/)

## Quick Start guide

### Step 1. Deploy Kubernetes operator using all in one config file

```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/main/deploy/all-in-one.yaml
```

### Step 2. Create Atlas Cluster

**1.** Create an Atlas API Key Secret

In order to work with the Atlas Operator you need to
provide [authentication information](https://docs.atlas.mongodb.com/configure-api-access)
to allow the Atlas Operator to communicate with Atlas API. Once you have generated a Public and Private key in Atlas,
you can create a Kuberentes Secret with:

```
kubectl create secret generic mongodb-atlas-operator-api-key \
         --from-literal="orgId=<the_atlas_organization_id>" \
         --from-literal="publicApiKey=<the_atlas_api_public_key>" \
         --from-literal="privateApiKey=<the_atlas_api_private_key>" \
         -n mongodb-atlas-system

kubectl label secret mongodb-atlas-operator-api-key atlas.mongodb.com/type=credentials -n mongodb-atlas-system
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

**3.** Create an `AtlasDeployment` Custom Resource.

The example below is a minimal configuration to create an M10 Atlas cluster in the AWS US East region. For a full list
of properties, check
`atlasclusters.atlas.mongodb.com` [CRD specification](config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml)):

```
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasDeployment
metadata:
  name: my-atlas-cluster
spec:
  projectRef:
    name: my-project
  clusterSpec:
    name: "Test-cluster"
    providerSettings:
      instanceSizeName: M10
      providerName: AWS
      regionName: US_EAST_1
EOF
```

**4.** Create a database user password Kubernetes Secret

```
kubectl create secret generic the-user-password --from-literal="password=P@@sword%"

kubectl label secret the-user-password atlas.mongodb.com/type=credentials
```

(note) To create X.509 user please see [this doc](docs/x509-user.md).

**5.** Create an `AtlasDatabaseUser` Custom Resource

In order to connect to an Atlas Cluster the database user needs to be created. `AtlasDatabaseUser` resource should
reference the password Kubernetes Secret created in the previous step.

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

Wait until the AtlasDatabaseUser resource gets to "ready" status (it will wait until the cluster is created that may
take around 10 minutes):

```
kubectl get atlasdatabaseusers my-database-user -o=jsonpath='{.status.conditions[?(@.type=="Ready")].status}'
True
```

### Step 3. Connect your application to the Atlas Cluster

The Atlas Operator will create a Kubernetes Secret with the information necessary to connect to the Atlas Cluster
created in the previous step. An application in the same Kubernetes Cluster can mount and use the Secret:

```
...
containers:
      - name: test-app
        env:
         - name: "CONNECTION_STRING"
           valueFrom:
             secretKeyRef:
               name: test-atlas-operator-project-test-cluster-theuser
               key: connectionStringStandardSrv

```

## Additional information or features

Operator support Third Party Integration.

- [Mongodb Atlas Operator sample](docs/project-integration.md)
- [Atlas documentation Atlas](https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings/)

### Step 4. Test Database as a Service (DBaaS) on OpenShift

The Atlas Operator is integrated with the [Red Hat Database-as-a-Service (DBaaS) Operator](https://github.com/RHEcosystemAppEng/dbaas-operator) which allows application developers to import database instances and connect to the databases through the [Service Binding Operator](https://github.com/redhat-developer/service-binding-operator). More information can be found [here](https://github.com/RHEcosystemAppEng/dbaas-operator#readme).

Note that both the DBaaS Operator and Atlas Operator should be installed through the [Operator Lifecyle Manager (OLM)](https://github.com/operator-framework/operator-lifecycle-manager).

**1.** Check DBaaS Registration

If the DBaaS Operator has been deployed in the OpenShift Cluster, the Atlas Operator automatically creates a cluster level [DBaaSProvider](https://github.com/RHEcosystemAppEng/dbaas-operator/blob/main/config/crd/bases/dbaas.redhat.com_dbaasproviders.yaml) custom resource (CR) object `mongodb-atlas-registration` to automatically register itself with the DBaaS Operator. 

```
apiVersion: dbaas.redhat.com/v1alpha1
kind: DBaaSProvider
metadata:
  labels:
    related-to: dbaas-operator
    type: dbaas-provider-registration
  name: mongodb-atlas-registration
spec:
  connectionKind: MongoDBAtlasConnection
  credentialFields:
  - displayName: Organization ID
    key: orgId
    required: true
    type: string
  - displayName: Public API Key
    key: publicApiKey
    required: true
    type: string
  - displayName: Private API Key
    key: privateApiKey
    required: true
    type: maskedstring
  inventoryKind: MongoDBAtlasInventory
  provider:
    displayDescription: Cloud-hosted MongoDB service on AWS, Azure and Google Cloud
    displayName: MongoDB Atlas Cloud Database Service
    icon:
      base64data: <icon image data>
      mediatype: image/png
    name: Red Hat DBaaS / MongoDB Atlas
```
If the Atlas Operator is undeployed with the OLM, the above registration CR gets cleaned up automatically.

**2.** Check MongoDBAtlasInventory Custom Resource

First an administrator creates a [DBaaSInventory](https://github.com/RHEcosystemAppEng/dbaas-operator/blob/main/config/crd/bases/dbaas.redhat.com_dbaasinventories.yaml) CR for MongoDB. The DBaaS Operator automatically creates a MongoDBAtlasInventory CR, and the Atlas Operator discovers the clusters and  instances, and sets the result in the CR status. 
Here is an example of MongoDBAtlasInventory CR.
```
apiVersion: dbaas.redhat.com/v1alpha1
kind: MongoDBAtlasInventory
metadata:
  name: dbaas-mytest
  namespace: openshift-operators
  ownerReferences:
  - apiVersion: dbaas.redhat.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: DBaaSInventory
    name: dbaas-mytest
    uid: 01f5a690-c640-462f-b6e8-ccb9db95df70
spec:
  credentialsRef:
    name: my-atlas-key
    namespace: openshift-operators
status:
  conditions:
  - lastTransitionTime: "2021-08-18T20:06:47Z"
    message: Spec sync OK
    reason: SyncOK
    status: "True"
    type: SpecSynced
  instances:
  - instanceID: 608df625aa94426b4169999
    instanceInfo:
      dbaas-cluster1.a00aa.Srv: mongodb+srv://dbaas-cluster1.a00aa.mongodb.net
      instanceSizeName: M0
      projectID: 608df5e652e1944293e11111
      projectName: Project 1
      providerName: TENANT
      regionName: US_EAST_1
    name: DBaaS-Cluster1
  - instanceID: 60807282b4ab8d3b3c123456
    instanceInfo:
       dbaas-cluster1.a00aa.Srv: mongodb+srv://test123.edbca.mongodb.net
      instanceSizeName: M10
      projectID: 6065e15b16c0731bf3a12333
      projectName: Project 2
      providerName: AWS
      regionName: US_EAST_1
    name: test
  - instanceID: 12345ffbc9a90e310e642482
    instanceInfo:
      dbaas-cluster1.a00aa.Srv: mongodb+srv://testcluster1.edbca.mongodb.net
      instanceSizeName: M0
      projectID: 6065e15b16c0731bf3a12333
      projectName: Project 2
      providerName: TENANT
      regionName: US_EAST_1
    name: DBCreatedInAtalas
  - instanceID: 60b7a72f4877d05881234567
    instanceInfo:
      dbaas-cluster1.a00aa.Srv: mongodb+srv://test.abcd9.mongodb.net
      instanceSizeName: M10
      projectID: 60b798fea37f9f09acc12345
      projectName: mytest
      providerName: AWS
      regionName: US_EAST_1
    name: test
```
**3.** Check MongoDBAtlasConnection Custom Resource

Now the application developer can create a [DBaaSConnection](https://github.com/RHEcosystemAppEng/dbaas-operator/blob/main/config/crd/bases/dbaas.redhat.com_dbaasconnections.yaml) CR for connection to the MongoDB database instance found, the DBaaS Operator automatically creates a MongoDBAtlasConnection CR. The Atlas Operator creates a database user in Atlas for the cluster with the default database `admin`. The Atlas Operator stores the db user credentials in a kubernetes secret, and the remaining connection information in a configmap, and then updates the MongoDBAtlasConnection CR status.

Here is an example of MongoDBAtlasConnection CR.
```
apiVersion: dbaas.redhat.com/v1alpha1
kind: MongoDBAtlasConnection
metadata:
  name: test-dbaas-connection
  namespace: test-namespace
  ownerReferences:
  - apiVersion: dbaas.redhat.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: DBaaSConnection
    name: test-dbaas-connection
    uid: 77193619-6ab1-43c9-acf2-a40c2cfe7703
spec:
  instanceID: 12345ffbc9a90e310e642482
  inventoryRef:
    name: dbaas-mytest
    namespace: openshift-operators
status:
  conditions:
  - lastTransitionTime: "2021-08-18T20:07:51Z"
    message: ""
    reason: Ready
    status: "True"
    type: ReadyForBinding
  connectionInfoRef:
    name: atlas-connection-cm-knp9z
  credentialsRef:
    name: atlas-db-user-5pc8b
```
The corresponding generated secret:
```
apiVersion: v1
data:
  password: cGFzczEyM3dAcmQ=
  username: ZGJVc2VyXzEwMQ==
kind: Secret
metadata:
  labels:
    managed-by: atlas-operator
    owner: test-dbaas-connection
    owner.kind: MongoDBAtlasConnection
    owner.namespace: test-namespace
  name: atlas-db-user-5pc8b
  namespace: test-namespace
  ownerReferences:
  - apiVersion: dbaas.redhat.com/v1alpha1
    blockOwnerDeletion: false
    controller: true
    kind: MongoDBAtlasConnection
    name: test-dbaas-connection
    uid: a50b06db-8fa1-45c9-9893-833a028dfccc
type: Opaque
```
The corresponding generated configmap:
```
apiVersion: v1
data:
  host: cluster0.ubajs.mongodb.net
  provider: Red Hat DBaaS / MongoDB Atlas
  srv: "true"
  type: mongodb
kind: ConfigMap
metadata:
   labels:
    managed-by: atlas-operator
    owner: test-dbaas-connection
    owner.kind: MongoDBAtlasConnection
    owner.namespace: test-namespace
    name: atlas-connection-cm-knp9z
  namespace: test-namespace
  ownerReferences:
  - apiVersion: dbaas.redhat.com/v1alpha1
    blockOwnerDeletion: false
    controller: true
    kind: MongoDBAtlasConnection
    name: test-dbaas-connection
    uid: a50b06db-8fa1-45c9-9893-833a028dfccc
```
## How to Contribute

Please file issues before filing PRs. For PRs to be accepted, contributors must sign
our [CLA](https://www.mongodb.com/legal/contributor-agreement).

Reviewers, please ensure that the CLA has been signed by referring
to [the contributors tool](https://contributors.corp.mongodb.com/) (internal link).
