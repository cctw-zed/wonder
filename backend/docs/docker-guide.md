# Docker Operations Guide

This document records the standard Docker workflows for the Wonder project.

## Build

```bash
docker compose build
```

Builds the application image and ancillary monitoring services (Prometheus, Grafana, ELK components).

## Start Stack

```bash
docker compose up -d
```

Starts Wonder, PostgreSQL, Prometheus, Grafana, Elasticsearch, Logstash, Kibana, and cAdvisor in detached mode.

## Stop Stack

```bash
docker compose down
```

Stops all containers and removes the default network while preserving named volumes.

## View Logs

```bash
docker compose logs -f
```

Streams logs from all services. To inspect a single service:

```bash
docker compose logs -f wonder
```

## Rebuild & Restart Application Only

```bash
docker compose build wonder
```

```bash
docker compose up -d wonder
```

Useful when only application code changes.

## Database Access

```bash
docker exec -it wonder-postgres psql -U wonder postgres
```

Opens a PostgreSQL shell for administrative tasks. To run a one-shot query as the development user:

```bash
docker exec -i wonder-postgres psql -U dev -d wonder_dev -c 'SELECT 1;'
```

## Monitoring URLs

- Wonder API: `http://localhost:8080`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`
- Kibana: `http://localhost:5601`
- Elasticsearch: `http://localhost:9200`
- cAdvisor: `http://localhost:8081`

Ensure Docker Desktop (or the selected engine) is running before executing these commands.
