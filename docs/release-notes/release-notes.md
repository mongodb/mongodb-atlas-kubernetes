# MongoDB Atlas Operator v1.2.0

## Atlas Operator 

* Updated to Go 1.18 #604

## AtlasProject Resource

* Added support for Private Endpoints backwards sync #603

## AtlasDeployment Resource

* Refactored the Advanced Deployment Handler #615 (#606)
* Changed autoScaling to a new struct according to Atlas API #592 (#588)
* Fixed diskSizeGB decreasing for normal deployments #634 (#611)
* Fixed panic when Atlas API returns an empty object #593 (#589)

*The images can be found in:*

https://quay.io/mongodb/mongodb-atlas-kubernetes-operator