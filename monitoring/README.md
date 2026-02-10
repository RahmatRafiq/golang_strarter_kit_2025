# Observability Stack (Prometheus + Grafana)

This directory contains configuration files for monitoring and observability using Prometheus and Grafana.

## Quick Start

### 1. Start Monitoring Stack

```bash
docker-compose -f docker-compose.monitoring.yml up -d
```

### 2. Start Application

```bash
go run main.go
# or
./your-app-binary
```

### 3. Access Services

- **Application**: http://localhost:9999
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: `admin`
- **Metrics Endpoint**: http://localhost:9999/metrics

## Available Metrics

### HTTP Metrics
- `http_requests_total` - Total number of HTTP requests (by method, path, status)
- `http_request_duration_seconds` - Request latency histogram (p50, p95, p99)
- `http_response_size_bytes` - Response size histogram
- `http_requests_in_flight` - Current number of requests being served

### Database Metrics
- `db_connection_pool_stats` - Connection pool statistics
  - `max_open_connections` - Maximum number of open connections
  - `open_connections` - Current open connections
  - `in_use` - Connections currently in use
  - `idle` - Idle connections
  - `wait_count` - Total number of connections waited for
  - `wait_duration_ms` - Total time waited for connections
  - `max_idle_closed` - Connections closed due to SetMaxIdleConns
  - `max_lifetime_closed` - Connections closed due to SetConnMaxLifetime

## Grafana Dashboards

The Grafana instance comes pre-configured with:
1. **Request Rate** - Requests per second
2. **Latency Percentiles** - p50, p95, p99 response times
3. **Error Rate** - 4xx and 5xx error percentages
4. **Database Connection Pool** - Connection utilization
5. **Status Code Distribution** - Breakdown by HTTP status codes

## Example Prometheus Queries

### Request Rate (req/s)
```promql
rate(http_requests_total[5m])
```

### p95 Latency
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Error Rate (%)
```promql
100 * (
  sum(rate(http_requests_total{status=~"5.."}[5m]))
  /
  sum(rate(http_requests_total[5m]))
)
```

### Database Connection Pool Usage (%)
```promql
100 * (
  db_connection_pool_stats{metric="in_use"}
  /
  db_connection_pool_stats{metric="max_open_connections"}
)
```

## Configuration Files

- `prometheus.yml` - Prometheus scrape configuration
- `grafana/datasources/prometheus.yml` - Grafana datasource config
- `grafana/dashboards/dashboard.yml` - Dashboard provisioning config

## Stopping Services

```bash
docker-compose -f docker-compose.monitoring.yml down
```

## Troubleshooting

### Prometheus can't scrape metrics
- Check if application is running on port 9999
- Verify `/metrics` endpoint is accessible: `curl http://localhost:9999/metrics`
- Check Prometheus targets page: http://localhost:9090/targets

### Grafana shows no data
- Verify Prometheus datasource is configured correctly
- Check if Prometheus is collecting metrics
- Ensure time range in Grafana matches when app was running
