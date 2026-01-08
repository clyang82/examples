# Realistic Performance Test Guide

This guide helps you run performance tests based on your actual use cases with wildcards and path-based resource naming.

## Your 7 Test Cases

The performance test simulates these real-world scenarios:

1. **Global Admin**: `user1` can access all resources on all clusters
2. **Multi-Cluster Admin**: `user2` can access all resources on cluster1 and cluster2 only
3. **Resource-Type Scoped**: `user3` can access all pods on cluster1 (all namespaces)
4. **Namespace Scoped**: `user4` can access pods on cluster1 in namespace1 only
5. **Specific Permission**: `user5` can create pods on cluster1 in namespace1
6. **Global Resource-Type**: `user6` can access pods on all clusters
7. **Group Membership**: `user7` has access via group1 membership

## Quick Start

Run the complete realistic test:

```bash
./run-flat-test.sh
```

This will:
1. Generate test data with your 7 test cases + bulk data
2. Start SpiceDB
3. Load the schema and relationships
4. Run performance benchmarks
5. Verify each test case
6. Generate results

## Schema Design

### Realistic Schema (Flat with Wildcards)

The `schema_flat.yaml` uses a **flat structure** with **path-based resource naming**:

```
resource:cluster/_wildcard_                                    → All clusters
resource:cluster/cluster1/namespace/_wildcard_                 → All namespaces in cluster1
resource:cluster/cluster1/namespace/namespace1/core/pods/_wildcard_ → All pods in ns1
resource:cluster/cluster1/namespace/namespace1/core/pods/pod1       → Specific pod
```

**Key features:**
- ✅ Single resource definition
- ✅ Wildcards for pattern matching
- ✅ Path-based naming for hierarchy representation
- ✅ No complex permission inheritance
- ✅ Direct relationships

### vs. Hierarchical Schema

Your current `schema_new.yaml` uses **separate definitions** for cluster, namespace, and resource:

```
cluster → namespace → resource (with permission inheritance via ->)
```

## Resource Naming Convention

### Path-based Format

```
resource:cluster/{clusterId}/namespace/{namespaceId}/{apiGroup}/{resourceType}/{resourceName}
```

**Examples:**
- `resource:cluster/cluster1/namespace/namespace1/core/pods/pod1`
- `resource:cluster/cluster1/namespace/namespace1/apps/deployments/nginx`
- `resource:cluster/cluster1/namespace/namespace1/core/services/api`

### Wildcard Usage

Use `_wildcard_` to match multiple resources:

```
cluster/_wildcard_                           → All clusters
cluster/cluster1/namespace/_wildcard_        → All namespaces in cluster1
cluster/cluster1/namespace/ns1/core/pods/_wildcard_ → All pods in ns1
```

## Performance Comparison

### Expected Results

| Aspect | Flat + Wildcards | Hierarchical Inheritance |
|--------|------------------|------------------------|
| **Permission Checks** | Fast (direct lookup) | Slower (traverses hierarchy) |
| **Wildcard Matching** | Native support | Requires multiple relationships |
| **Bulk Assignments** | Easy (one wildcard relation) | Easy (assign at parent level) |
| **Relationship Count** | Moderate | Lower |
| **Latency** | 2-5ms | 5-15ms |
| **Debugging** | Simple (direct paths) | Complex (inheritance chains) |
| **Flexibility** | Very high | Moderate |

## Understanding Test Results

### Benchmark Metrics

```
Test Case                                Avg        P95        P99     Throughput
------------------------------------------------------------------------
Case 1: Wildcard all clusters           2.5ms      4.1ms      6.2ms    400.0/s
Case 2: Multiple specific clusters      2.3ms      3.8ms      5.9ms    435.0/s
Case 3: Pods in all namespaces          2.7ms      4.3ms      6.5ms    370.0/s
```

**What to look for:**
- **Avg Latency**: Typical permission check time
- **P95/P99**: Worst-case scenarios (important for SLAs)
- **Throughput**: How many checks per second
- **Success Count**: Should match expected permissions

## Test Data Structure

### Generated Data

For 100 clusters, 100 namespaces, 100 pods:

```
Relationships:
  - 7 test case relationships (wildcards)
  - 5 cluster-level admins (wildcards)
  - 10 namespace-level editors (wildcards)
  - 1,000,000 pod resources (specific)
  - 1,000,000 viewer assignments

Total: ~1,000,022 relationships
```

### Comparison with Hierarchical

**Flat + Wildcards:**
```
resource:cluster/_wildcard_#admin@user:user1
→ 1 relationship covers all resources
```

**Hierarchical:**
```
cluster:cluster0#admin@user:user1
cluster:cluster1#admin@user:user1
...
cluster:cluster99#admin@user:user1
→ 100 relationships + inheritance chain
```

## Manual Testing

### Check Specific Permission

```bash
zed permission check \
  resource:cluster/cluster1/namespace/namespace1/core/pods/pod1 \
  get \
  user:user1 \
  --endpoint localhost:50051 \
  --token foobar \
  --insecure
```

### Lookup All Resources for User

```bash
zed permission lookup-resources \
  resource \
  get \
  user:user1 \
  --endpoint localhost:50051 \
  --token foobar \
  --insecure
```

### Verify Group Membership

```bash
# Check user7's access via group1
zed permission check \
  resource:cluster/cluster1/namespace/namespace1/core/pods/pod1 \
  get \
  user:user7 \
  --endpoint localhost:50051 \
  --token foobar \
  --insecure
```

## Customizing Tests

### Adjust Test Scale

Edit `run-flat-test.sh`:

```bash
# Small test (faster)
CLUSTERS=10
NAMESPACES=10
PODS=10

# Medium test
CLUSTERS=50
NAMESPACES=50
PODS=50

# Large test (closer to production)
CLUSTERS=100
NAMESPACES=100
PODS=100

# Very large test (stress test)
CLUSTERS=200
NAMESPACES=200
PODS=200
```

### Add More Test Cases

Edit `generate-flat-test-data.go`:

```go
// Add custom test case
fmt.Fprintf(file, "  // Case 8: Custom test case\n")
fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/namespace1/apps/deployments/_wildcard_#admin@user:user8\n\n")
```

## Comparing Schemas

### Run Both Tests

```bash
# Run realistic (flat + wildcards) test
./run-flat-test.sh

# Run hierarchical test
./run-performance-test.sh
```

### Compare Results

```bash
./compare-results.sh \
  performance-results/flat_20240101_120000.txt \
  performance-results/hierarchical_20240101_120000.txt
```

## Key Decision Factors

### Choose Flat + Wildcards if:

✅ You need **flexible pattern matching** (e.g., all pods, specific resource types)
✅ **Low latency** is critical (< 5ms)
✅ You have **path-based resource naming** (like Kubernetes)
✅ **Wildcard support** is important
✅ You want **simple debugging** (direct paths)
✅ Your access patterns use **resource type filtering** (all pods, all deployments)

### Choose Hierarchical if:

✅ You have **strict organizational hierarchy** (cluster owns namespace owns resource)
✅ **Permission inheritance** is natural (cluster admin → namespace admin → resource admin)
✅ You want **fewer total relationships**
✅ **Moderate latency** is acceptable (< 15ms)
✅ You rarely need **cross-cutting access** (e.g., all pods on all clusters)

### Real-World Recommendation

Based on your 7 test cases, **Flat + Wildcards** is likely better because:

1. ✅ **Case 1, 6**: Require wildcard matching across all clusters (hard in hierarchical)
2. ✅ **Case 3**: Resource-type scoped access (all pods) is natural with wildcards
3. ✅ **Case 2**: Multi-cluster access easier with wildcards vs. multiple inheritance paths
4. ✅ **Case 4, 5, 7**: Specific namespace access works in both, but simpler in flat

**However, Hierarchical might be better if:**
- Most users are scoped to single clusters (hierarchy matches user model)
- You have very deep nesting (cluster → namespace → workload → container)
- Permission changes propagate down hierarchy frequently

## Performance Optimization Tips

### 1. Use Wildcards Efficiently

**Good:**
```
resource:cluster/_wildcard_/namespace/_wildcard_/core/pods/_wildcard_#viewer@user:pod_viewer
```
→ One relationship for all pods

**Bad:**
```
resource:cluster/cluster0/namespace/namespace0/core/pods/pod0#viewer@user:pod_viewer
resource:cluster/cluster0/namespace/namespace0/core/pods/pod1#viewer@user:pod_viewer
...
```
→ One relationship per pod

### 2. Index Common Patterns

If most checks are for specific resource types (pods, deployments), consider:
- Pre-computing common wildcard matches
- Using SpiceDB caching
- Optimizing your resource naming convention

### 3. Monitor P99 Latency

Even with fast average latency, high P99 can cause issues:
- Set alerts on P99 > 10ms
- Investigate slow queries with `--show-debug-traces`

## Troubleshooting

### Wildcard Not Matching

If wildcards aren't working, check:

```bash
# Verify relationship exists
zed relationship read \
  resource:cluster/_wildcard_ \
  --endpoint localhost:50051 \
  --token foobar \
  --insecure

# Check permission with debug traces
zed permission check \
  resource:cluster/cluster1/namespace/namespace1/core/pods/pod1 \
  get \
  user:user1 \
  --endpoint localhost:50051 \
  --token foobar \
  --insecure \
  --show-debug-traces
```

### Slow Performance

If latency is high:

1. **Check relationship count**: `zed relationship read | wc -l`
2. **Use persistent datastore**: Memory datastore has limits
3. **Enable caching**: Configure SpiceDB cache
4. **Optimize queries**: Use more specific wildcards

### Memory Issues

For large datasets (> 1M resources):

```bash
# Use PostgreSQL instead of memory
docker run -d --name spicedb-postgres \
  -e POSTGRES_PASSWORD=secret \
  -p 5432:5432 \
  postgres:15

spicedb migrate head \
  --datastore-engine postgres \
  --datastore-conn-uri "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"

spicedb serve \
  --grpc-preshared-key foobar \
  --datastore-engine postgres \
  --datastore-conn-uri "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
```

## Next Steps

1. **Run the test**: `./run-flat-test.sh`
2. **Review results**: Check `performance-results/`
3. **Analyze metrics**: Focus on P95/P99 for your SLA
4. **Compare approaches**: Run both flat and hierarchical tests
5. **Make decision**: Choose schema based on your needs

## Files Overview

- `schema_flat.yaml` - Flat schema with wildcard support
- `generate-flat-test-data.go` - Generates your 7 test cases + bulk data
- `benchmark-flat.go` - Benchmarks all 7 test cases
- `run-flat-test.sh` - Automated test runner
- `relationships_flat.yaml` - Generated test data

## Support

For SpiceDB-specific questions:
- Docs: https://authzed.com/docs
- Discord: https://authzed.com/discord
- GitHub: https://github.com/authzed/spicedb
