# Automated secret provisioning with Secrets Store CSI Driver

The Atlas Kubernetes Operator consumes credentials from Kubernetes Secrets, continuously waiting for changes in them.

With this, the Operator is ready to support both manually or automatically deployed secrets. Among the automated secrets, one popular solution is to provision the secrets directly off the company's or team's Vault service of choice. One Open Source tool allowing for this automation is [Secrets Store CSI Driver](https://secrets-store-csi-driver.sigs.k8s.io/). Note, however, [the main usage of the CSI driver is to mount secrets as local files in pods, but it also supports the option to expose Kubernetes Secrets our of those mounts so long as some pod is mounting them](https://secrets-store-csi-driver.sigs.k8s.io/topics/sync-as-kubernetes-secret).

In this directory we have a couple of examples on how you would integrate the Operator Atlas credentials and DB User Password with a Hashicorp Vault.

For more details, please check out the documentation at [Secrets Store CSI Driver](https://secrets-store-csi-driver.sigs.k8s.io/) directly.

Here are the pointers to the Secrets Store provider plugins:

- [Hashicorp Vault](https://github.com/hashicorp/vault-csi-provider).
- [AWS plugin provider](https://github.com/aws/secrets-store-csi-driver-provider-aws) Secrets Manager or Parameters Store.
- [Azure Key Vault](https://azure.github.io/secrets-store-csi-driver-provider-azure/docs/).
- [Google Cloud Secret Manager](https://github.com/GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp).
