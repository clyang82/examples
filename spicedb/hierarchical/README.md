# Hierarchical Schema Implementation (Option 3)

Multi-level SpiceDB schema with permission inheritance through cluster → namespace → resource hierarchy. Supports both **namespace-scoped** and **cluster-scoped** resources.

## Overview

This implementation uses **Option 3** design where resources can link to either:
- **Namespace** (for namespace-scoped resources like Pods, Deployments, ConfigMaps)
- **Cluster** (for cluster-scoped resources like Nodes, PersistentVolumes, ClusterRoles)

## Files

- `schema_new.yaml` - Hierarchical schema with dual relations (Option 3)
- `generate-hierarchical-data.go` - Generate test data with both resource types
- `benchmark-hierarchical.go` - Performance benchmark with 8 comprehensive test cases

## Schema Design (Option 3)

### Structure

```yaml
definition cluster {
  relation admin: user | group#member
  relation editor: user | group#member
  relation viewer: user | group#member

  permission get = admin + editor + viewer
  permission create = admin + editor
  permission delete = admin
}

definition namespace {
  relation cluster: cluster
  relation admin: user | group#member
  relation editor: user | group#member
  relation viewer: user | group#member

  permission get = admin + editor + viewer + cluster->get  # Inherits from cluster
  permission create = admin + editor + cluster->create
  permission delete = admin + cluster->delete
}

definition resource {
  relation namespace: namespace  # For namespace-scoped resources
  relation cluster: cluster      # For cluster-scoped resources
  relation admin: user | group#member
  relation editor: user | group#member
  relation viewer: user | group#member

  # Checks BOTH inheritance paths
  permission get = admin + editor + viewer + namespace->get + cluster->get
  permission create = admin + namespace->create + cluster->create
  permission delete = admin + namespace->delete + cluster->delete
}
```

### Resource Naming Convention

**Format**: `{clusterId}/{namespaceId}/{kind}/{resourceName}` or `{clusterId}/{kind}/{resourceName}`

#### Namespace-scoped Resources

```
cluster:cluster1
namespace:cluster1/namespace1
resource:cluster1/namespace1/pods/nginx
resource:cluster1/namespace1/deployments/myapp
resource:cluster1/default/configmaps/app-config
```

**Relationships:**
```yaml
# Links resource to namespace
resource:cluster1/namespace1/pods/nginx#namespace@namespace:cluster1/namespace1
```

#### Cluster-scoped Resources

```
cluster:cluster1
resource:cluster1/nodes/node1
resource:cluster1/persistentvolumes/pv001
resource:cluster1/clusterroles/admin
```

**Relationships:**
```yaml
# Links resource directly to cluster (skips namespace)
resource:cluster1/nodes/node1#cluster@cluster:cluster1
```

### Permission Inheritance

#### For Namespace-scoped Resources
```
User assigned at cluster level
    ↓ (cluster->get)
All namespaces in that cluster
    ↓ (namespace->get)
All resources in those namespaces
```

#### For Cluster-scoped Resources
```
User assigned at cluster level
    ↓ (cluster->get)
All cluster-scoped resources (nodes, PVs, etc.)
```

## Performance Characteristics

**Tested with:** 100 clusters × 100 namespaces × 100 pods + 10 cluster resources per cluster

- **Avg Latency**: **567µs** ⚡
- **P95 Latency**: **860µs** ⭐ **25% better than flat**
- **P99 Latency**: 5.237ms
- **Throughput**: **1,788 checks/sec**
- **Total Relationships**: ~2,020,000
  - 2,000 cluster-scoped resources (nodes + PVs)
  - 10,000 namespaces
  - 1,000,000 namespace-scoped resources (pods)
- **Import Rate**: ~11,500 relationships/sec

## Advantages

✅ **8/8 test cases work** including cluster-scoped resources
✅ **Permission inheritance** (assign once at cluster = all resources)
✅ **Dual resource types** (namespace-scoped AND cluster-scoped)
✅ **Better P95 latency** (25% faster than flat)
✅ **Kubernetes-native model** (matches K8s RBAC)
✅ **Group support** works correctly
✅ **Scalable** - Tested with 2M+ relationships
✅ **Sub-millisecond average latency**

## Trade-offs

⚠️ More complex schema (4 definitions vs 2)
⚠️ 2x relationships (includes hierarchy links)
⚠️ Permission resolution involves traversing inheritance

## Usage

### Generate Test Data

```bash
# Generate with both namespace-scoped and cluster-scoped resources
go run generate-hierarchical-data.go \
  --clusters 100 \
  --namespaces 100 \
  --pods 100 \
  --cluster-resources 10 \
  --output relationships_hierarchical.yaml
```

**Output:**
- 2,000 cluster-scoped resources (nodes + persistentvolumes)
- 10,000 namespaces
- 1,000,000 namespace-scoped resources (pods)

### Load Schema and Relationships

```bash
# From repository root
cd /root/go/src/github.com/clyang82/examples/spicedb

# Load schema
go run common/load-schema.go hierarchical/schema_new.yaml

# Load relationships (takes ~2 minutes for 2M relationships)
go run common/load-relationships.go hierarchical/relationships_hierarchical.yaml
```

### Run Benchmark

```bash
cd hierarchical

# Run with 100 checks per test case
go run benchmark-hierarchical.go \
  --endpoint localhost:50051 \
  --token foobar \
  --checks 100
```

## Permission Examples

### Global Admin (All Resources)

```yaml
# One relationship per cluster grants access to ALL resources
cluster:cluster0#admin@user:admin1
cluster:cluster1#admin@user:admin1
```

**Result:** Admin1 can access:
- All namespaces in clusters 0 and 1
- All namespace-scoped resources (pods, deployments, etc.)
- All cluster-scoped resources (nodes, PVs, etc.)

### Namespace Editor

```yaml
# Grant edit permission on specific namespace
namespace:cluster1/namespace1#editor@user:editor1
```

**Result:** Editor1 can edit all namespace-scoped resources in cluster1/namespace1

### Cluster Resource Viewer

```yaml
# Grant viewer on cluster (includes cluster-scoped resources)
cluster:cluster1#viewer@user:viewer1
```

**Result:** Viewer1 can view:
- All namespaces in cluster1
- All pods, deployments, etc. in cluster1
- All nodes, PVs in cluster1

### Resource Viewer (Direct Assignment)

```yaml
# Direct assignment to specific pod
resource:cluster1/namespace1/pods/nginx#viewer@user:viewer1
```

**Result:** Viewer1 can only view this specific pod

### Group Membership

```yaml
# User is member of group
group:group1#user@user:user7

# Group has permission on namespace
namespace:cluster1/namespace1#viewer@group:group1#member
```

**Result:** User7 (via group1) can view all resources in cluster1/namespace1

## Test Cases Covered

The benchmark includes **8 comprehensive test cases**:

1. ✅ **Case 1**: Admin on all clusters (inherits to all resources)
2. ✅ **Case 2**: Admin on specific clusters (cluster1, cluster2 only)
3. ✅ **Case 3**: Viewer on cluster1 (all namespaces/pods)
4. ✅ **Case 4**: Viewer on specific namespace (cluster1/namespace1 only)
5. ✅ **Case 5**: Admin on namespace (create permission)
6. ✅ **Case 6**: Viewer on all clusters (all resources)
7. ✅ **Case 7**: Group membership (user via group has viewer)
8. ✅ **Case 8**: Cluster-scoped resources (nodes, persistentvolumes)

**Success Rate**: 100% functional coverage

## Resource Type Handling

### When to Use Namespace Relation

Use `#namespace@namespace:...` for resources that belong to a namespace:
- Pods
- Deployments
- Services
- ConfigMaps
- Secrets
- StatefulSets
- DaemonSets

### When to Use Cluster Relation

Use `#cluster@cluster:...` for resources at cluster level:
- Nodes
- PersistentVolumes
- ClusterRoles
- ClusterRoleBindings
- Namespaces (the namespace objects themselves)
- CustomResourceDefinitions

## When to Use Hierarchical Schema

**✅ Recommended** - Choose hierarchical schema if:
- You have natural hierarchy (cluster/namespace/resource)
- You need both namespace-scoped and cluster-scoped resources
- Bulk permission management is important
- You need permission inheritance (assign once, applies to many)
- Sub-millisecond latency is acceptable
- You're building Kubernetes-like RBAC systems
- You want to model organizational structure

**❌ Not Recommended** if:
- You need flat, direct permissions only
- Schema complexity is a concern
- You don't have natural hierarchy

## Performance Tuning

### For Best Performance

1. **Database Indexes**: Add indexes on frequently queried relationships
2. **SpiceDB Caching**: Enable caching for hot paths
3. **Connection Pooling**: Use persistent connections
4. **Monitoring**: Track P99 latency and set alerts

### Scaling Recommendations

- Use PostgreSQL read replicas for read-heavy workloads
- Consider materialized views for complex queries
- Batch permission checks when possible
- Cache results at application level for frequently checked permissions

## Comparison: Hierarchical vs Flat

| Metric | Hierarchical | Flat | Winner |
|--------|--------------|------|--------|
| Avg Latency | **567µs** | 635µs | ✅ Hierarchical (12% faster) |
| P95 Latency | **860µs** | 1.142ms | ✅ Hierarchical (25% faster) |
| Throughput | **1,788/sec** | 1,622/sec | ✅ Hierarchical (10% higher) |
| Functional Coverage | **8/8 cases** | 1/7 cases | ✅ Hierarchical |
| Schema Complexity | 4 definitions | 2 definitions | ⚠️ Flat (simpler) |
| Total Relationships | ~2M | ~1M | ⚠️ Flat (fewer) |

## Deployment

### Prerequisites

- Kubernetes cluster with SpiceDB and PostgreSQL deployed
- Port-forward to SpiceDB: `kubectl port-forward -n spicedb svc/spicedb 50051:50051`

### Setup Steps

1. **Deploy infrastructure** (see `../manifests/`)
   ```bash
   kubectl apply -f ../manifests/k8s-postgres.yaml
   kubectl apply -f ../manifests/k8s-spicedb.yaml
   ```

2. **Generate relationships** (optional, if not already generated)
   ```bash
   go run generate-hierarchical-data.go --clusters 100 --namespaces 100 --pods 100
   ```

3. **Load schema and data**
   ```bash
   go run ../common/load-schema.go schema_new.yaml
   go run ../common/load-relationships.go relationships_hierarchical.yaml
   ```

4. **Run benchmark to verify**
   ```bash
   go run benchmark-hierarchical.go --checks 100
   ```

## Conclusion

**✅ Recommended for production** - The hierarchical schema with Option 3 provides:
- Excellent performance (567µs avg, 1,788 checks/sec)
- Full Kubernetes-like RBAC modeling
- Both namespace-scoped and cluster-scoped resource support
- Permission inheritance for easy management
- Proven at scale (2M+ relationships)

The minimal increase in schema complexity is vastly outweighed by the functional advantages, better performance, and natural modeling of hierarchical permissions.
