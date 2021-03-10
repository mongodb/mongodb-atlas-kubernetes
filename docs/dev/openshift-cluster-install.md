## Openshift cluster installation

Open the link (you need to have an account registered in RedHat and added to "mongoDB" organization):
https://cloud.redhat.com/openshift/install/aws/installer-provisioned

1. Ensure you have an AWS account configured in `~/.aws/credentials`
1. Download and unpack the MacOS installer
1. Copy the `scripts/openshift/install-config.yaml` to some `<temp_directory>`. Set the following:
   * `pullSecret: '<..>'` (copy the content of the Pull Secrets from the link above)
   * `sshKey: | \n` (public ssh Key)
1. Run `./openshift-install create cluster --dir=<temp_directory>`.
  * `--dir=` points to the directory where the `install-config.yaml` is located
  * the installer will verify permissions and will show the ones that are missing - it's necessary to give them to your AWS account
1. Wait for ~40 minutes

Some notes on configuration of the cluster:
* it's not possible to have less than 3 replicas for control plane and 2 replicas for worker nodes 
  (see https://docs.openshift.com/container-platform/4.7/installing/installing_aws/installing-aws-customizations.html?extIdCarryOver=true&intcmp=7013a000002CtetAAC&sc_cid=701f2000001OH7iAAG#installation-configuration-parameters_installing-aws-customizations)
* by default Openshift uses `m5.xlarge` for controlPlane nodes and `m5.large` for worker nodes. `m5.xlarge` uses 16Gb for memory
  and this is the smallest memory allowed for each controlplane instance. In our development we use `t3.xlarge` which provides the same
  memory but costs cheaper.
* removing the existing cluster can be done by calling `./openshift-install destroy cluster` (not sure if this needs the SSH keys)

The log of the last installation:

```
➜  ./openshift-install create cluster --dir=/Users/alisovenko/temp-openshift
INFO Credentials loaded from default AWS environment variables
INFO Consuming Install Config from target directory
WARNING Missing permissions to fetch Quotas and therefore will skip checking them: failed to load limits for servicequotas: failed to list serviceqquotas for ec2: AccessDeniedException: User: arn:aws:iam::268558157000:user/anton.lisovenko is not authorized to perform: servicequotas:ListServiceQuotas, make sure you have `servicequotas:ListAWSDefaultServiceQuotas` permission available to the user.
INFO Creating infrastructure resources...
INFO Waiting up to 20m0s for the Kubernetes API at https://api.atlas.operator.mongokubernetes.com:6443...
INFO API v1.20.0+bd9e442 up
INFO Waiting up to 30m0s for bootstrapping to complete...
INFO Destroying the bootstrap resources...
INFO Waiting up to 40m0s for the cluster at https://api.atlas.operator.mongokubernetes.com:6443 to initialize...
INFO Waiting up to 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/Users/alisovenko/workspace/mongodb-atlas-kubernetes/scripts/openshift/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.atlas.operator.mongokubernetes.com
INFO Login to the console with user: "kubeadmin", and password: "*****"
INFO Time elapsed: 35m55s
```