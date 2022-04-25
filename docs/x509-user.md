# Generate X.509 Client Certificates

## [Prerequisite] Generate an X.509 certificate

### A. Using [cert-manager](https://cert-manager.io/docs/)

- Install the [cert-manager](https://cert-manager.io/docs/installation/)
  ```
  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
  ```

- Create an Issuer ([example](./helpers/certs/issuer.yaml))
  ```
  kubectl apply -f docs/helpers/certs/issuer.yaml
  ```

- Create a Certificate ([example](./helpers/certs/cert.yaml))
  ```
  kubectl apply -f docs/helpers/certs/cert.yaml`
  ```

### B. Using a custom [script](../scripts/create_x509.go)

- Run the [script](../scripts/create_x509.go):
  ```
  go run scripts/create_x509.go --path=tmp/x509/
  ```

- Put cert into the secret:
  ```
  kubectl create secret generic my-x509-cert --from-file=./tmp/x509/cert.pem
  kubectl label secret my-x509-cert atlas.mongodb.com/type=credentials
  ```

## Enable X.509 for project and create a user 

### Create a Project:
```yaml
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject     
metadata:
  name: my-project      
spec:                 
  name: Test Project
  projectIpAccessList:     
    - ipAddress: "192.0.2.15"
      comment: "IP address for Application Server A"
    - ipAddress: "203.0.113.0/24"
      comment: "CIDR block for Application Server B - D"
  x509CertRef:
    name: my-x509-cert
EOF
```

### Create a User:
```yaml
cat <<EOF | kubectl apply -f -
apiVersion: atlas.mongodb.com/v1
kind: AtlasDatabaseUser
metadata:
  name: my-database-user
spec:
  username: CN=my-x509-authenticated-user,OU=organizationalunit,O=organization
  databaseName: "\$external"
  x509Type: "CUSTOMER"
  roles:
    - roleName: "readWriteAnyDatabase"
      databaseName: "admin"
  projectRef:
    name: my-project
EOF
```