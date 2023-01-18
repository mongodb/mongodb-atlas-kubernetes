# MongoDB Atlas Operator v1.5.0

## Operator Changes

* Fix connection secret creation (#774)
* Fix mininum version of Openshift (#761)
* add operator class to only reconcile resources with mongodb.com/atlas-operator-class=<operator-class> annotation (#842)

## AtlasProject Resource

* Add Atlas Teams (#767)
* Make sure PEs are always added to status (#773)
* Fix the InstanceSize must match issue #777 (#782)

## AtlasDeployment Resource

* Add serverless PE support (#779)
* Convert OplogMinRetentionHours field properly (#778)

*The images can be found in:*
https://quay.io/mongodb/mongodb-atlas-kubernetes-operator
