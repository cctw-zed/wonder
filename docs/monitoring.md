# Monitoring Stack Guide

This document outlines how to run the local observability stack for the Wonder project.

## Components

- **Prometheus** (`prom/prometheus:v2.53.0`): Scrapes Wonder's `/metrics` endpoint and cAdvisor for container metrics.
- **Grafana** (`grafana/grafana:11.2.0`): Pre-provisioned with a Prometheus datasource for dashboarding and alert rules.
- **Elasticsearch** (`docker.elastic.co/elasticsearch/elasticsearch:8.14.1`): Stores structured application logs.
- **Logstash** (`docker.elastic.co/logstash/logstash:8.14.1`): Ingests container logs via GELF and forwards them to Elasticsearch.
- **Kibana** (`docker.elastic.co/kibana/kibana:8.14.1`): Visualizes log data and builds searches/alerts.
- **cAdvisor** (`gcr.io/cadvisor/cadvisor:v0.49.1`): Exposes container resource usage for Prometheus scraping.

## Running the Stack

```bash
# build application image and start entire stack
docker compose build
docker compose up -d
```

Services are exposed on the following ports:

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)
- Elasticsearch: `http://localhost:9200`
- Logstash GELF input: UDP `12201`
- Kibana: `http://localhost:5601`
- cAdvisor: `http://localhost:8081`
- Wonder API: `http://localhost:8080`

> **Note:** On macOS, cAdvisor host volume mounts may require additional permissions. If the container fails, comment out the `cadvisor` service in `docker-compose.yaml` or adjust the mounts.

## Metrics

Wonder now exposes Prometheus-compatible metrics at `/metrics` on port `8080`. The middleware tracks request counts and latency histograms per HTTP method and route. Prometheus scrapes the `wonder` job every 15 seconds using the configuration in `monitoring/prometheus/prometheus.yml`.

To verify metrics:

1. Open Prometheus at `http://localhost:9090` and run queries such as:
   - `wonder_http_requests_total`
   - `rate(wonder_http_request_duration_seconds_sum[1m])`
2. In Grafana, import dashboards for Gin/Go services or build custom panels using the provisioned Prometheus datasource.

## Logs

Container logs from the Wonder service are shipped through the Docker GELF logging driver to Logstash, which structures the records and stores them in Elasticsearch (index pattern `wonder-logs-*`).

Steps to inspect logs:

1. Start Kibana (`http://localhost:5601`).
2. Configure an index pattern for `wonder-logs-*`.
3. Use Discover to search logs or create visualizations.

## Alerting & Next Steps

- Define Prometheus alert rules for latency, error rate, and database connectivity in `monitoring/prometheus/prometheus.yml`.
- Use Grafana alerting to route notifications to Slack or email.
- Extend Logstash pipelines for additional services or enrich log fields with trace IDs.
- Consider adding Filebeat or Fluent Bit if richer log collection is needed beyond GELF.
