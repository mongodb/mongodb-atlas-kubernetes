# Automated secret provisioning with External Secrets

The Atlas Kubernetes Operator consumes credentials from Kubernetes Secrets, continuously waiting for changes in them.

With this, the Operator is ready to support both manually or automatically deployed secrets. Among the automated secrets, one popular solution is to provision the secrets directly off the company's or team's Vault service of choice. One Open Source tool allowing for this automation is [External Secrets](https://external-secrets.io/latest/).

In this directory we have a couple of examples on how you would integrate the Operator Atlas credentials and DB User Password with a Hashicorp Vault.

For more details, please check out the documentation at [External Secrets](https://external-secrets.io/latest/) directly.

Here are some pointers to popular External Secrets providers:

- [Hashicorp Vault](https://external-secrets.io/latest/provider/hashicorp-vault/).
- [AWS Secrets Manager](https://external-secrets.io/latest/provider/aws-secrets-manager/) or [Parameters Store](https://external-secrets.io/latest/provider/aws-parameter-store/).
- [Azure Key Vault](https://external-secrets.io/latest/provider/azure-key-vault/).
- [Google Cloud Secret Manager](https://external-secrets.io/latest/provider/google-secrets-manager/).
