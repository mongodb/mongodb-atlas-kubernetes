# Create a PrivateLink for GCP

## I. Create a Private Endpoint Service
```yaml
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: Test Atlas Operator Project
  privateEndpoints:
  - provider: "GCP"
    region: "us-east1"
EOF
```

## II. Setup GCP Side of connection

- Use `EDIT` Atlas UI Button and follow a few steps to get a similar script:

  ```
  #!/bin/bash
  gcloud config set project atlasoperator

  for i in {0..5}
  do
    gcloud compute addresses create user-private-endpoint-ip-$i --region=us-east1 --subnet=user-test-subnet
  done

  for i in {0..5}
  do
    if [ $(gcloud compute addresses describe user-private-endpoint-ip-$i --region=us-east1 --format="value(status)") != "RESERVED" ]; then
      echo "user-private-endpoint-ip-$i is not RESERVED";
      exit 1;
    fi
  done

  for i in {0..5}
  do
    gcloud compute forwarding-rules create user-private-endpoint-$i --region=us-east1 --network=user-test-vpc --address=user-private-endpoint-ip-$i --target-service-attachment=projects/p-long-id/regions/us-east1/serviceAttachments/long-id-$i
  done

  if [ $(gcloud compute forwarding-rules list --regions=us-east1 --format="csv[no-heading](name)" --filter="name:user-private-endpoint" | wc -l) -gt 6 ]; then
    echo "Project has too many forwarding rules that match prefix user-private-endpoint. Either delete the competing resources or choose another endpoint prefix."
    exit 2;
  fi

  gcloud compute forwarding-rules list --regions=us-east1 --format="json(IPAddress,name)" --filter="name:user-private-endpoint" > atlasEndpoints-user-private-endpoint.json
  ```

- Run the scipt `sh setup_psk.sh`
- Run a couple command to format the output for the operator:
  ```bash
  yq e -P atlasEndpoints-user-private-endpoint.json > atlasEndpoints-user-private-endpoint.yaml
  awk 'sub("name","endpointName")sub("IPAddress","ipAddress")' atlasEndpoints-user-private-endpoint.yaml
  ```
  Expected output:
  ```
  - ipAddress: 10.0.0.00
    endpointName: user-private-endpoint-0
  - ipAddress: 10.0.0.01
    endpointName: user-private-endpoint-1
  - ipAddress: 10.0.0.02
    endpointName: user-private-endpoint-2
  - ipAddress: 10.0.0.03
    endpointName: user-private-endpoint-3
  - ipAddress: 10.0.0.04
    endpointName: user-private-endpoint-4
  - ipAddress: 10.0.0.05
    endpointName: user-private-endpoint-5
  ```

## III. Create the Private Endpoint Inteface
```yaml
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: Test Atlas Operator Project
  privateEndpoints:
  - provider: "GCP"
    region: "us-east1"
    gcpProjectId: "atlasoperator"
    endpointGroupName: "user-test-vpc"
    endpoints:
    - ipAddress: 10.0.0.00
      endpointName: user-private-endpoint-0
    - ipAddress: 10.0.0.01
      endpointName: user-private-endpoint-1
    - ipAddress: 10.0.0.02
      endpointName: user-private-endpoint-2
    - ipAddress: 10.0.0.03
      endpointName: user-private-endpoint-3
    - ipAddress: 10.0.0.04
      endpointName: user-private-endpoint-4
    - ipAddress: 10.0.0.05
      endpointName: user-private-endpoint-5
EOF
```