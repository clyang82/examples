# Kubernetes Deployment

Kubernetes manifests for deploying SpiceDB with PostgreSQL backend.

## Files

- `k8s-postgres.yaml` - PostgreSQL 15 deployment
- `k8s-spicedb.yaml` - SpiceDB deployment with migrations

## Architecture

```
┌─────────────────┐
│   SpiceDB Pod   │
│  (Port: 50051)  │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│ PostgreSQL Pod  │
│  (Port: 5432)   │
└─────────────────┘
```

## Deployment Steps

### 1. Deploy PostgreSQL

```bash
kubectl apply -f k8s-postgres.yaml
kubectl wait --for=condition=ready pod -l app=postgres -n spicedb --timeout=120s
```

**Resources:**
- Namespace: `spicedb`
- Deployment: `postgres`
- Service: `postgres:5432`
- Storage: 5Gi PVC
- Image: `postgres:15`

**Configuration:**
- Database: `spicedb`
- User: `spicedb`
- Password: `spicedb123`

### 2. Deploy SpiceDB

```bash
kubectl apply -f k8s-spicedb.yaml
kubectl wait --for=condition=complete job/spicedb-migrate -n spicedb --timeout=60s
kubectl wait --for=condition=ready pod -l app=spicedb -n spicedb --timeout=60s
```

**Resources:**
- Deployment: `spicedb`
- Service: `spicedb:50051`
- Job: `spicedb-migrate` (runs schema migrations)
- Image: `quay.io/authzed/spicedb:v1.35.0`

**Configuration:**
- Token: `foobar`
- Datastore: PostgreSQL
- HTTP: Disabled
- gRPC: Port 50051

### 3. Access SpiceDB

```bash
# Port-forward to local machine
kubectl port-forward -n spicedb svc/spicedb 50051:50051

# Test connection
grpcurl -plaintext localhost:50051 list
```

## Resource Specifications

### PostgreSQL

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

### SpiceDB

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "2Gi"
    cpu: "1000m"
```

## Health Checks

Both deployments include:
- **Readiness probes**: Check service availability
- **Liveness probes**: Detect and restart unhealthy pods

## Persistence

PostgreSQL uses a PersistentVolumeClaim:
```yaml
accessModes:
  - ReadWriteOnce
storage: 5Gi
```

**Note**: Data persists across pod restarts

## Scaling

### PostgreSQL
```bash
# Not recommended - use replication instead
kubectl scale deployment postgres -n spicedb --replicas=1
```

### SpiceDB
```bash
# Scale for higher throughput
kubectl scale deployment spicedb -n spicedb --replicas=3
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n spicedb
```

### View Logs

```bash
# PostgreSQL logs
kubectl logs -n spicedb -l app=postgres

# SpiceDB logs
kubectl logs -n spicedb -l app=spicedb

# Migration job logs
kubectl logs -n spicedb -l job-name=spicedb-migrate
```

### Check Migration Status

```bash
kubectl get jobs -n spicedb
kubectl describe job spicedb-migrate -n spicedb
```

### Restart Pods

```bash
# Restart SpiceDB
kubectl delete pod -l app=spicedb -n spicedb

# Restart PostgreSQL (data persists)
kubectl delete pod -l app=postgres -n spicedb
```

### Clean Up

```bash
# Delete everything
kubectl delete namespace spicedb

# Delete just SpiceDB (keep PostgreSQL)
kubectl delete deployment spicedb -n spicedb
kubectl delete job spicedb-migrate -n spicedb
```

## Production Considerations

For production deployments, consider:

1. **Security**
   - Use TLS for gRPC
   - Store credentials in Secrets
   - Use strong passwords
   - Enable network policies

2. **High Availability**
   - Deploy PostgreSQL with replication
   - Run multiple SpiceDB replicas
   - Use pod disruption budgets

3. **Performance**
   - Increase resource limits
   - Use dedicated nodes
   - Enable connection pooling
   - Configure PostgreSQL tuning

4. **Monitoring**
   - Add Prometheus metrics
   - Set up alerts
   - Monitor query performance
   - Track error rates

5. **Backup**
   - Regular PostgreSQL backups
   - Test restore procedures
   - Use volume snapshots

## Configuration Options

### SpiceDB Environment Variables

```yaml
SPICEDB_GRPC_PRESHARED_KEY: "your-token"
SPICEDB_DATASTORE_ENGINE: "postgres"
SPICEDB_DATASTORE_CONN_URI: "postgres://..."
SPICEDB_HTTP_ENABLED: "false"
```

### PostgreSQL Environment Variables

```yaml
POSTGRES_DB: "spicedb"
POSTGRES_USER: "spicedb"
POSTGRES_PASSWORD: "your-password"
```

## Testing

After deployment:

```bash
# Load schema
go run ../common/load-schema.go ../hierarchical/schema_new.yaml

# Import relationships
go run ../common/load-relationships.go relationships.yaml

# Run benchmark
go run ../hierarchical/benchmark-hierarchical.go
```

## Network

Services are ClusterIP by default. For external access:

```bash
# Port-forward (development)
kubectl port-forward -n spicedb svc/spicedb 50051:50051

# LoadBalancer (production)
# Edit k8s-spicedb.yaml and change Service type to LoadBalancer
```

## Storage

Default storage class is used. To specify:

```yaml
storageClassName: "your-storage-class"
```

For high performance, use SSD-backed storage.
