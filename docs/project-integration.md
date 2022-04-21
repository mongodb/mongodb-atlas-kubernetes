# Project with Third Party Integration

### Sample

Operator support 3rd Party Integration with following [services](https://www.mongodb.com/docs/atlas/reference/api/third-party-integration-settings-create/)

For example DATADOG requare (according documentation):
- type = DATADOG
- apiKey = Your API Key.
- region = Indicates which API URL to use, either US or EU. Datadog will use US by default.


```
apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: TestAtlasOperatorProjectIntegration3
  connectionSecretRef:
    name: my-atlas-key
  projectIpAccessList:
    - ipAddress: "0.0.0.0/1"
      comment: "Everyone has access. For the test purpose only."
    - ipAddress: "128.0.0.0/1"
      comment: "Everyone has access. For the test purpose only."
  integrations:
    - type: "DATADOG"
      apiKeyRef:
        name: key-name
        namespace: key-namespace
      region: "US"
```
