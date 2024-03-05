# MongoDB Atlas Operator

[![MongoDB Atlas Operator](https://github.com/mongodb/mongodb-atlas-kubernetes/workflows/Test/badge.svg)](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/test.yml?query=branch%3Amain)
[![MongoDB Atlas Go Client](https://img.shields.io/badge/Powered%20by%20-go--client--mongodb--atlas-%2313AA52)](https://github.com/mongodb/go-client-mongodb-atlas)

The MongoDB Atlas Operator provides a native integration between the Kubernetes orchestration platform and MongoDB Atlas
— the only multi-cloud document database service that gives you the versatility you need to build sophisticated and
resilient applications that can adapt to changing customer demands and market trends.

The full documentation for the Operator can be found [here](https://docs.atlas.mongodb.com/atlas-operator/)

## Quick Start guide

### Prerequisites to Install using kubectl

Before you install the MongoDB Atlas Operator using `kubectl`, you must:

1. Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).
2. Have Kubernetes installed, i.e. [Kind](https://kind.sigs.k8s.io/)
3. Choose the [version of the operator](https://github.com/mongodb/mongodb-atlas-kubernetes/releases) you want to run

### Step 1. Deploy Kubernetes operator using all in one config file

```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/<operator_version>/deploy/all-in-one.yaml
```

### Step 2. Create Atlas Deployment

**1.** Create an Atlas API Key Secret

In order to work with the Atlas Operator you need to
provide [authentication information](https://docs.atlas.mongodb.com/configure-api-access)
to allow the Atlas Operator to communicate with Atlas API. Once you have generated a Public and Private key in Atlas,
you can create a Kuberentes Secret with:

```
kubectl create secret generic mongodb-atlas-operator-api-key \
         --from-literal='orgId=<the_atlas_organization_id>' \
         --from-literal='publicApiKey=<the_atlas_api_public_key>' \
         --from-literal='privateApiKey=<the_atlas_api_private_key>' \
         -n mongodb-atlas-system

kubectl label secret mongodb-atlas-operator-api-key atlas.mongodb.com/type=credentials -n mongodb-atlas-system
```

**2.** Create an `AtlasProject` Custom Resource

The `AtlasProject` CustomResource represents Atlas Projects in our Kubernetes cluster. You need to specify
`projectIpAccessList` with the IP addresses or CIDR blocks of any hosts that will connect to the Atlas Deployment.

```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/<operator_version>/config/samples/atlas_v1_atlasproject.yaml
```

**3.** Create an `AtlasDeployment` Custom Resource.

The example below is a configuration to create an M10 Atlas deployment in the AWS US East region. For a full list
of properties, check
`atlasdeployments.atlas.mongodb.com` [CRD specification](config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml)):

```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/<operator_version>/config/samples/atlas_v1_atlasdeployment.yaml
```

**4.** Create a database user password Kubernetes Secret

```
kubectl create secret generic the-user-password --from-literal='password=P@@sword%'

kubectl label secret the-user-password atlas.mongodb.com/type=credentials
```

(note) To create X.509 user please see [this doc](docs/x509-user.md).

**5.** Create an `AtlasDatabaseUser` Custom Resource

In order to connect to an Atlas Deployment the database user needs to be created. `AtlasDatabaseUser` resource should
reference the password Kubernetes Secret created in the previous step.

```
kubectl apply -f https://raw.githubusercontent.com/mongodb/mongodb-atlas-kubernetes/<operator_version>/config/samples/atlas_v1_atlasdatabaseuser.yaml
```

**6.** Wait for the `AtlasDatabaseUser` Custom Resource to be ready

Wait until the AtlasDatabaseUser resource gets to "ready" status (it will wait until the deployment is created that may
take around 10 minutes):

```
kubectl get atlasdatabaseusers my-database-user -o=jsonpath='{.status.conditions[?(@.type=="Ready")].status}'
True
```

### Step 3. Connect your application to the Atlas Deployment

The Atlas Operator will create a Kubernetes Secret with the information necessary to connect to the Atlas Deployment
created in the previous step. An application in the same Kubernetes Cluster can mount and use the Secret:

```
...
containers:
- name: test-app
  env:
    - name: "CONNECTION_STRING"
      valueFrom:
        secretKeyRef:
          name: test-atlas-operator-project-test-deployment-theuser
          key: connectionStringStandardSrv

```

## Additional information or features

In certain cases you can modify the default operator behaviour via [annotations](docs/annotations.md).

Operator support Third Party Integration.

- [Mongodb Atlas Operator sample](docs/project-integration.md)
- [Atlas documentation Atlas](https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings/)

## How to Contribute

Please file issues before filing PRs. For PRs to be accepted, contributors must sign
our [CLA](https://www.mongodb.com/legal/contributor-agreement).

Reviewers, please ensure that the CLA has been signed by referring
to [the contributors tool](https://contributors.corp.mongodb.com/) (internal link).
