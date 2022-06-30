# Setup Prometheus and Grafana

## Steps

### 1. Create `prometheus.yaml` config file

This is an example:
```
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "Test Atlas Operator Project-mongo-metrics"
    scrape_interval: 10s
    metrics_path: /metrics
    scheme : https
    basic_auth:
      username: <user>
      password: <password>
    http_sd_configs:
      - url: https://cloud.mongodb.com/prometheus/v1.0/groups/{group-ID}/discovery
        refresh_interval: 60s
        basic_auth:
          username: <user>
          password: <password>
```

You can find the URL in atlasproject status via a command:
```
kubectl get atlasproject -o json
```
Look for `prometheusDiscoveryURL` field.


### 2. Run prometheus with you config

Docker example:
```
docker run \
    -p 9090:9090 \
    -v /path/to/local/prometheus.yaml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```

Go to http://localhost:9090/ to check if it is working

### 3. Run grafana

Docker example:
```
docker run --name grafana -p 3000:3000 grafana/grafana
```

Go to http://localhost:3000/. Default username/password is `admin`.

For Docker Data Source `localhost` won't work and should be switched to `host.docker.internal`.

### 4. Create a new Data Source

- Go to Configuration > Data sources (aka http://localhost:3000/datasources)
- Add a new data source
- Set `Name`
- Set `URL` to you prometheus url (ex. http://host.docker.internal:9090/)
- Set `Basic auth` with credentials from you prometheus config
- Save & test

### 5. Create a new Dashboard

- Go to Dashboards (aka http://localhost:3000/dashboards)
- Import a new dashboard using [grafana example](./grafana/sample_dashboard.json)
- Set `Datasource` to the Data Source `Name` you've set in the previous step
- Metrics should appear
