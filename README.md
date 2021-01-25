# MongoDB Atlas Operator
Welcome to MongoDB Atlas Operator - the Operator allowing to manage Atlas Clusters from Kubernetes.

> Current Status: *pre-alpha*. We are currently working on the initial set of features that will give users the opportunity to
provision Atlas projects, clusters and database users using Kubernetes Specifications and bind connection information 
into the applications deployed to Kubernetes.   

## Installation / Upgrade

### Using Kubernetes config files

```
kubectl apply -f  <link to all_in_one>
```

## Create Atlas Cluster

In order to work with Atlas you'll need to provide the authentication information that would allow the Atlas Operator
communicate with Atlas API. You need to create the secret first:

```
kubectl create secret generic my-atlas-key \
         --from-literal="orgId=<the_atlas_organization_id>" \
         --from-literal="publicApiKey=<the_atlas_api_public_key>" \
         --from-literal="privateApiKey=<the_atlas_api_private_key>" \
         -n <your_namespace>
```

