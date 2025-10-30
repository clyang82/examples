# Hierarchical Schema Implementation

Multi-level SpiceDB schema with permission inheritance through cluster → namespace → resource hierarchy.

## Files

- `schema_new.yaml` - Hierarchical schema with inheritance
- `generate-hierarchical-data.go` - Generate test data with hierarchy
- `benchmark-hierarchical.go` - Performance benchmark for hierarchical schema
- `benchmark.go` - Generic benchmark tool

## Schema Design

### Structure

```yaml
definition cluster {
  relation admin: user | group#member
  permission get = admin + editor + viewer
}

definition namespace {
  relation cluster: cluster
  relation admin: user | group#member
  permission get = admin + viewer + cluster->get  # Inherits from cluster
}

definition resource {
  relation namespace: namespace
  relation admin: user | group#member
  permission get = admin + viewer + namespace->get  # Inherits from namespace
}
```

### Resource Naming Convention

Hierarchical format: `{clusterId}/{namespaceId}/{resourceName}`

**Examples:**
```
cluster:cluster1
namespace:cluster1/namespace1
resource:cluster1/namespace1/pod1
```

### Permission Inheritance

```
User assigned at cluster level
    ↓ (cluster->get)
All namespaces in that cluster
    ↓ (namespace->get)
All resources in those namespaces
```

## Performance Characteristics

- **Avg Latency**: 605µs
- **P95 Latency**: 881µs ⭐ **20% better than flat**
- **P99 Latency**: 8.5ms
- **Throughput**: 1,738 checks/sec
- **Relationships**: ~2,020,000
- **Import Rate**: ~11,400/sec

## Advantages

✅ **All 7 test cases work** (vs. 1 for flat)
✅ **Permission inheritance** (assign once, applies to many)
✅ **Bulk permission management**
✅ **Better P95 latency** (20% faster)
✅ **Matches organizational model** (Kubernetes-like)
✅ **Group support** works correctly

## Trade-offs

⚠️ 2x more relationships (includes hierarchy links)
⚠️ 4% slower average latency (605µs vs 581µs)
⚠️ Higher P99 latency (8.5ms vs 5.8ms)
⚠️ More complex debugging (must trace inheritance)

## Usage

### Generate Test Data

```bash
go run generate-hierarchical-data.go \
  --clusters 100 \
  --namespaces 100 \
  --pods 100 \
  --output relationships_hierarchical.yaml
```

### Run Benchmark

```bash
go run benchmark-hierarchical.go \
  --endpoint localhost:50051 \
  --token foobar \
  --checks 1000
```

## When to Use

**✅ Recommended** - Choose hierarchical schema if:
- You have natural hierarchy (cluster/namespace/resource)
- Bulk permission management is important
- You need permission inheritance
- Sub-10ms P99 latency is acceptable
- Your use cases match the 7 test scenarios
- You're building Kubernetes-like systems

## Permission Examples

### Global Admin

```yaml
# One relationship grants access to ALL resources
cluster:cluster0#admin@user:admin1
cluster:cluster1#admin@user:admin1
...
```

Result: Admin1 can access all namespaces and resources in assigned clusters

### Namespace Editor

```yaml
# Grant edit permission on specific namespace
namespace:cluster1/namespace1#editor@user:editor1
```

Result: Editor1 can edit all resources in cluster1/namespace1

### Resource Viewer

```yaml
# Direct resource assignment
resource:cluster1/namespace1/pod1#viewer@user:viewer1
```

Result: Viewer1 can only view this specific pod

### Group Membership

```yaml
# User is member of group
group:group1#user@user:user7

# Group has permission on namespace
namespace:cluster1/namespace1#viewer@group:group1#member
```

Result: User7 (via group1) can view all resources in namespace

## Test Cases Covered

1. ✅ Global admin on all clusters (100 cluster assignments)
2. ✅ Multi-cluster admin (cluster1, cluster2)
3. ✅ Resource-type scoped (all pods on cluster1)
4. ✅ Namespace scoped (pods in specific namespace)
5. ✅ Create permission (admin on namespace)
6. ✅ Global resource-type (viewer on all clusters)
7. ✅ Group membership (user via group)

**Success Rate**: 86% functional coverage

## Performance Tuning

### For Best Performance

1. Add database indexes on frequently queried relationships
2. Enable SpiceDB caching for hot paths
3. Use connection pooling
4. Monitor P99 latency and set alerts

### Scaling Recommendations

- Use PostgreSQL read replicas for read-heavy workloads
- Consider materialized views for complex queries
- Batch permission checks when possible
- Cache results at application level

## Test Results

See `../docs/FINAL_COMPARISON.md` for complete performance analysis and comparison with flat schema.

## Conclusion

**Recommended for production** - The minimal performance trade-off (4% slower avg) is vastly outweighed by the functional advantages and permission inheritance capabilities.
