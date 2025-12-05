# MongoDB Atlas Operator

[![MongoDB Atlas Operator](https://github.com/mongodb/mongodb-atlas-kubernetes/workflows/Test/badge.svg)](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/test.yml?query=branch%3Amain)

The MongoDB Atlas Operator provides a native integration between the Kubernetes orchestration platform and MongoDB Atlas
— the only multi-cloud document database service that gives you the versatility you need to build sophisticated and
resilient applications that can adapt to changing customer demands and market trends.

The full documentation for the Operator can be found [here](https://docs.atlas.mongodb.com/atlas-operator/)

## Getting Started

### Supported features

* Create and configure an Atlas Project, or connect to an existing one.
* Deploy, manage, scale, and tear down Atlas clusters.
* Support for Atlas serverless instances.
* Create and edit database users.
* Manage IP Access Lists, network peering and private endpoints.
* Configure and control Atlas’s fully managed cloud backup.
* Configure federated authentication for your Atlas organization
* Integrate Atlas monitoring with Prometheus.

... and more.

To view the list of custom resources and their respective schemas, visit our [reference](https://www.mongodb.com/docs/atlas/operator/stable/custom-resources/) 
documentation. See the [Quickstart](https://www.mongodb.com/docs/atlas/operator/stable/ak8so-quick-start/) to get started
with Atlas Kubernetes Operator.

If you just want to check the released images signatures, please [check these instructions](docs/dev/signed-images.md).

## How to Contribute

Please file issues before filing PRs. For PRs to be accepted, contributors must sign
our [CLA](https://www.mongodb.com/legal/contributor-agreement).

Reviewers, please ensure that the CLA has been signed by referring
to [the contributors tool](https://contributors.corp.mongodb.com/) (internal link).
