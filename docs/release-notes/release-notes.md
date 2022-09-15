# MongoDB Atlas Operator v1.3.0

## AtlasProject Resource

* Add network peering feature #620
* Add cloud provider access role feature #645
* Add encryption at rest #674

## AtlasDeployment Resource

* Fix deployment CR deletion if token invalid #666 (#421)
* Prevent changing instanceSize and diskGB if autoscaling is enabled #672 (#648, #649)
* Fix error message for Delete method #664 
* Add test for atlasdeployments with keep annotation #612

*The images can be found in:*

https://quay.io/mongodb/mongodb-atlas-kubernetes-operator
