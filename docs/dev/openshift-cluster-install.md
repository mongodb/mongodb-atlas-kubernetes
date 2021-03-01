## Openshift cluster installation

Open the link (you need to have an account registered in RedHat and added to "mongoDB" organization):
https://cloud.redhat.com/openshift/install/aws/installer-provisioned

1. Ensure you have an AWS account configured in `~/.aws/credentials`
1. Download and unpack the MacOS installer
1. Run `./openshift-install create cluster`. (*TODO this default configuration results in quite a big cluster, this needs 
   to be changed to a smaller values somehow*)
  * choose the zone that has enough VPCs
  * specify the Pull Secrets (can be copied from the link above)
  * the installer will verify permissions and will show the ones that are missing - it's necessary to give them to your AWS account
1. Wait for ~40 minutes

The log of the last installation:

```
➜  openshift-install-mac ./openshift-install create cluster
? SSH Public Key /Users/alisovenko/.ssh/id_aws_rsa.pub
? Platform aws
INFO Credentials loaded from default AWS environment variables
? Region eu-west-3
? Base Domain mongokubernetes.com
? Cluster Name atlas.operator.openshift
? Pull Secret [? for help] **********
WARNING Missing permissions to fetch Quotas and therefore will skip checking them: failed to load limits for servicequotas: failed to list default serviceqquotas for ec2: AccessDeniedException: User: arn:aws:iam::268558157000:user/anton.lisovenko is not authorized to perform: servicequotas:ListAWSDefaultServiceQuotas, make sure you have `servicequotas:ListAWSDefaultServiceQuotas` permission available to the user.
INFO Creating infrastructure resources...
INFO Waiting up to 20m0s for the Kubernetes API at https://api.atlas.operator.openshift.mongokubernetes.com:6443...
INFO API v1.20.0+bd9e442 up
INFO Waiting up to 30m0s for bootstrapping to complete...
INFO Destroying the bootstrap resources...
INFO Waiting up to 40m0s for the cluster at https://api.atlas.operator.openshift.mongokubernetes.com:6443 to initialize...
INFO Waiting up to 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=/Users/alisovenko/Downloads/Soft/openshift-install-mac/auth/kubeconfig'
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.atlas.operator.openshift.mongokubernetes.com
INFO Login to the console with user: "kubeadmin", and password: "(erased)"
INFO Time elapsed: 36m53s
```