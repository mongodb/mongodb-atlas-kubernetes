# SRE Runbook: Resource Stuck in Reconciliation

## Problem: Resource stuck in reconciliation
This problem showcases the issue with AtlasProject resource not being ready. It can be applied to every AKO resource

### Symptoms:
- The resource is not ready.
- High error rate metric.
  
To monitor the error rate, you can create a query to calculate the reconciliation error rate for the `AtlasProject` controller as a percentage over the last minute. This metric helps in identifying and monitoring the health and stability of the `AtlasProject` controller. A high or rising error percentage indicates issues in the reconciliation process.

#### Example Query:
To calculate the error rate, use the following Prometheus query:
```prometheus
100 * rate(controller_runtime_reconcile_errors_total{controller="AtlasProject"}[1m]) / rate(controller_runtime_reconcile_total{controller="AtlasProject"}[1m])
```

### Status:
Check the resource status condition for further details:
```yaml
status:
  conditions:
    - type: Ready
      status: "False"
      reason: ....
```

### Action Items:
1. **Verify Resource Status:**
   - Check the status condition message for more detailed information.
   - If the `AtlasProject` is not ready, proceed with the next troubleshooting steps.
  
2. **Check Connection Secret:**
   - Ensure the connection secret referenced by `spec.connectionSecretRef.name` is correctly labeled with `atlas.mongodb.com/type=credentials`.

3. **Investigate Logs:**
   - Review logs for the `AtlasProject` controller for any potential errors or failed reconciliation attempts.

### Additional Resources:
- [AtlasProject resource](https://www.mongodb.com/docs/atlas/operator/upcoming/atlasproject-custom-resource/)
