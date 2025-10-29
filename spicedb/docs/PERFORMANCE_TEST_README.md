# SpiceDB Performance Test Guide

This directory contains a complete performance testing suite to compare two SpiceDB schemas:
- **schema_flat.yaml** - Flat structure (resources only)
- **schema_new.yaml** - Hierarchical structure (cluster → namespace → resource)

## Prerequisites

1. **SpiceDB** installed and in your PATH
   ```bash
   # Install SpiceDB
   brew install authzed/tap/spicedb authzed/tap/zed
   # OR download from: https://github.com/authzed/spicedb/releases
   ```

2. **Go** 1.19 or later
   ```bash
   go version
   ```

3. **Required Go dependencies**
   ```bash
   go get github.com/authzed/authzed-go/v1
   go get github.com/authzed/authzed-go/proto/authzed/api/v1
   go get github.com/authzed/grpcutil
   go get google.golang.org/grpc
   ```

## Quick Start

Run the complete automated test suite:

```bash
./run-performance-test.sh
```

This will:
1. Generate test data for both schemas
2. Start SpiceDB instances
3. Load schemas and relationships
4. Run benchmarks
5. Generate comparison reports

## Manual Testing

If you prefer step-by-step control:

### Step 1: Generate Test Data

```bash
# For hierarchical schema
go run generate-test-data.go \
    --schema hierarchical \
    --clusters 100 \
    --namespaces 100 \
    --resources 100 \
    --output relationships_hierarchical.yaml

# For flat schema
go run generate-test-data.go \
    --schema flat \
    --clusters 100 \
    --namespaces 100 \
    --resources 100 \
    --output relationships_flat.yaml
```

### Step 2: Start SpiceDB

```bash
spicedb serve \
    --grpc-preshared-key foobar \
    --http-enabled=false \
    --datastore-engine=memory \
    --grpc-addr=":50051"
```

### Step 3: Load Schema

```bash
# For hierarchical schema
zed schema write schema_new.yaml \
    --endpoint localhost:50051 \
    --token foobar \
    --insecure

# For flat schema
zed schema write schema_flat.yaml \
    --endpoint localhost:50051 \
    --token foobar \
    --insecure
```

### Step 4: Import Relationships

```bash
# For hierarchical schema
zed relationship import relationships_hierarchical.yaml \
    --endpoint localhost:50051 \
    --token foobar \
    --insecure

# For flat schema
zed relationship import relationships_flat.yaml \
    --endpoint localhost:50051 \
    --token foobar \
    --insecure
```

### Step 5: Run Benchmarks

```bash
# For hierarchical schema
go run benchmark.go \
    --endpoint localhost:50051 \
    --token foobar \
    --schema hierarchical \
    --checks 1000

# For flat schema
go run benchmark.go \
    --endpoint localhost:50051 \
    --token foobar \
    --schema flat \
    --checks 1000
```

## Test Configuration

### Default Settings
- **Clusters**: 100
- **Namespaces**: 100 (distributed across clusters)
- **Resources per namespace**: 100
- **Total resources**: 10,000
- **Users**: 10 admins, 100 editors, 1,000 viewers

### Customizing Test Parameters

Modify the data generation:

```bash
go run generate-test-data.go \
    --clusters 50 \
    --namespaces 200 \
    --resources 50 \
    --admins 5 \
    --editors 50 \
    --viewers 500 \
    --schema hierarchical \
    --output custom_relationships.yaml
```

Adjust benchmark parameters:

```bash
go run benchmark.go \
    --checks 5000 \
    --concurrency 20 \
    --schema hierarchical
```

## Understanding the Results

### Hierarchical Schema Tests

The benchmark runs three types of permission checks:

1. **Cluster-level check** - Tests inherited permissions from cluster
   - User with cluster admin checks resource permission
   - Tests deep hierarchy traversal (cluster → namespace → resource)

2. **Namespace-level check** - Tests inherited permissions from namespace
   - User with namespace editor checks resource permission
   - Tests mid-level hierarchy traversal (namespace → resource)

3. **Resource-level check** - Tests direct resource permissions
   - User with resource viewer checks resource permission
   - Tests direct permission lookup (no traversal)

### Flat Schema Tests

The benchmark runs:

1. **Direct resource check** - Tests flat permission model
   - User checks direct resource permission
   - No hierarchy traversal

### Key Metrics

- **Avg Latency**: Average time per permission check
- **P50/P95/P99**: 50th, 95th, 99th percentile latencies
- **Throughput**: Checks per second
- **Errors**: Number of failed checks

## Expected Results Analysis

### Hierarchical Schema (schema_new.yaml)

**Advantages:**
- Fewer total relationships to manage
- Efficient bulk permission assignment (assign once at cluster/namespace level)
- Better organizational model for multi-tenant systems

**Trade-offs:**
- Slower permission checks due to hierarchy traversal
- P99 latency may be higher for deep hierarchies
- More complex to debug permission issues

### Flat Schema (schema_flat.yaml)

**Advantages:**
- Faster permission checks (direct lookup)
- More predictable latency (lower P99)
- Simpler debugging

**Trade-offs:**
- More relationships to manage (3x per resource)
- Harder to manage bulk permissions
- Less flexible organizational model

## Sample Results Format

```
Test: Cluster-level check
  Total Duration: 2.5s
  Avg Latency:    2.5ms
  P50 Latency:    2.3ms
  P95 Latency:    4.1ms
  P99 Latency:    6.2ms
  Throughput:     400.0 checks/sec
  Errors:         0
```

## Interpreting Results

### When to use Hierarchical Schema:
- You have clear organizational hierarchy (cluster/namespace/resource)
- You manage permissions at bulk level frequently
- You need flexible permission inheritance
- Total relationship count is a concern
- Latency requirements are moderate (< 10ms acceptable)

### When to use Flat Schema:
- You need minimal latency (< 5ms required)
- Permissions are managed per-resource
- Simple permission model is preferred
- Debugging simplicity is important
- Relationship count is not a concern

## Troubleshooting

### SpiceDB won't start
```bash
# Check if port is in use
lsof -i :50051

# Kill existing process
pkill spicedb
```

### Permission checks fail
```bash
# Verify schema is loaded
zed schema read --endpoint localhost:50051 --token foobar --insecure

# Check relationships
zed relationship read \
    --endpoint localhost:50051 \
    --token foobar \
    --insecure \
    resource:c0-ns0-r0
```

### Out of memory
```bash
# Use persistent datastore instead of memory
spicedb serve \
    --grpc-preshared-key foobar \
    --datastore-engine postgres \
    --datastore-conn-uri "postgres://user:pass@localhost/spicedb"
```

## Files Overview

- `performance-test-design.md` - Detailed test design and methodology
- `generate-test-data.go` - Data generation script
- `benchmark.go` - Performance benchmark tool
- `run-performance-test.sh` - Automated test runner
- `relationships_hierarchical.yaml` - Generated hierarchical test data
- `relationships_flat.yaml` - Generated flat test data
- `performance-results/` - Benchmark results directory

## Advanced Testing

### Testing with Persistent Storage

For more realistic testing, use PostgreSQL:

```bash
# Start PostgreSQL
docker run -d --name spicedb-postgres \
    -e POSTGRES_PASSWORD=secret \
    -p 5432:5432 \
    postgres:15

# Run SpiceDB with PostgreSQL
spicedb migrate head \
    --datastore-engine postgres \
    --datastore-conn-uri "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"

spicedb serve \
    --grpc-preshared-key foobar \
    --datastore-engine postgres \
    --datastore-conn-uri "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
```

### Stress Testing

Run longer benchmarks:

```bash
go run benchmark.go \
    --checks 10000 \
    --concurrency 50 \
    --schema hierarchical
```

### Profiling

Enable Go profiling:

```bash
go run benchmark.go \
    --checks 1000 \
    --cpuprofile cpu.prof \
    --memprofile mem.prof
```

## Next Steps

1. Run the automated test: `./run-performance-test.sh`
2. Review results in `performance-results/`
3. Adjust parameters based on your use case
4. Make schema decision based on your requirements

## Support

- SpiceDB Docs: https://authzed.com/docs
- SpiceDB Discord: https://authzed.com/discord
- GitHub Issues: https://github.com/authzed/spicedb/issues
