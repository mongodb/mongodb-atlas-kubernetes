
# MongoDB Atlas Operator v0.9.0

## Atlas Operator


## AtlasProject Resource
  * [3rd Party Integration](https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings/) are supported `spec.integrations`
  * [GCP Private Endpoints](https://www.mongodb.com/docs/atlas/reference/api/private-endpoints/) are now supported
  * [Maintenance Windows](https://www.mongodb.com/docs/atlas/reference/api/maintenance-windows/) are now supported

## AtlasDeployment Resource
* Changes
  * [Serverless instances](https://www.mongodb.com/docs/atlas/reference/api/serverless-instances/) are supported via the new `spec.serverlessSpec` field.

  
# MongoDB Atlas Operator v1.1.1

## Atlas Operator 

## AtlasDeployment Resource
* Changes
  * Change AtlasDeploymentSpec.advancedDeploymentSpec.replicationSpecs.regionConfigs.autoScaling to a new struct according to Atlas API. [Issue #588](https://github.com/mongodb/mongodb-atlas-kubernetes/issues/588)

*The images can be found in:*

https://quay.io/repository/mongodb/mongodb-atlas-operator