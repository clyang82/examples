# Schema Structure Comparison

## Hierarchical Schema (schema_new.yaml)

### Data Structure
```
cluster:cluster0 (admin: user:admin0)
├── namespace:ns0 (editor: user:editor0)
│   ├── resource:c0-ns0-r0 (viewer: user:viewer0)
│   ├── resource:c0-ns0-r1 (viewer: user:viewer1)
│   └── ... (100 resources)
├── namespace:ns1 (editor: user:editor1)
│   └── ... (100 resources)
└── ... (100 namespaces total across all clusters)

cluster:cluster1 (admin: user:admin1)
└── ...

... (100 clusters)
```

### Relationship Count
For 100 clusters, 100 namespaces, 100 resources per namespace:

| Relationship Type | Count | Example |
|------------------|-------|---------|
| Cluster admins | 10 | `cluster:cluster0#admin@user:admin0` |
| Namespace → Cluster | 100 | `namespace:ns0#cluster@cluster:cluster0` |
| Namespace editors | 100 | `namespace:ns0#editor@user:editor0` |
| Resource → Namespace | 10,000 | `resource:c0-ns0-r0#namespace@namespace:ns0` |
| Resource viewers | 10,000 | `resource:c0-ns0-r0#viewer@user:viewer0` |
| **Total** | **~20,110** | |

### Permission Evaluation Example

User `admin0` wants to `get` resource `c0-ns0-r0`:

```
1. Check: resource:c0-ns0-r0#admin@user:admin0? ❌ NO
2. Check: resource:c0-ns0-r0#editor@user:admin0? ❌ NO
3. Check: resource:c0-ns0-r0#viewer@user:admin0? ❌ NO
4. Check via namespace: resource:c0-ns0-r0#namespace->get@user:admin0?
   4a. Follow: resource:c0-ns0-r0#namespace → namespace:ns0
   4b. Check: namespace:ns0#get@user:admin0?
       4b1. Check: namespace:ns0#admin@user:admin0? ❌ NO
       4b2. Check: namespace:ns0#editor@user:admin0? ❌ NO
       4b3. Check: namespace:ns0#viewer@user:admin0? ❌ NO
       4b4. Check via cluster: namespace:ns0#cluster->get@user:admin0?
            4b4a. Follow: namespace:ns0#cluster → cluster:cluster0
            4b4b. Check: cluster:cluster0#get@user:admin0?
                  → Check: cluster:cluster0#admin@user:admin0? ✅ YES
5. Result: ✅ ALLOWED (via cluster inheritance)
```

**Traversal depth:** 3 levels (resource → namespace → cluster)

---

## Flat Schema (schema_flat.yaml)

### Data Structure
```
resource:res0
├── admin: user:admin0
├── editor: user:editor0
└── viewer: user:viewer0

resource:res1
├── admin: user:admin1
├── editor: user:editor1
└── viewer: user:viewer1

... (10,000 resources)
```

### Relationship Count
For 10,000 resources with admin, editor, and viewer each:

| Relationship Type | Count | Example |
|------------------|-------|---------|
| Resource admins | 10,000 | `resource:res0#admin@user:admin0` |
| Resource editors | 10,000 | `resource:res0#editor@user:editor0` |
| Resource viewers | 10,000 | `resource:res0#viewer@user:viewer0` |
| **Total** | **30,000** | |

### Permission Evaluation Example

User `admin0` wants to `get` resource `res0`:

```
1. Check: resource:res0#admin@user:admin0? ✅ YES
2. Result: ✅ ALLOWED (direct relationship)
```

**Traversal depth:** 1 level (resource only)

---

## Key Differences

| Aspect | Hierarchical | Flat |
|--------|-------------|------|
| **Relationships** | ~20,110 | 30,000 |
| **Traversal depth** | Up to 3 levels | 1 level |
| **Permission check complexity** | High (tree traversal) | Low (direct lookup) |
| **Bulk assignment** | Easy (assign at cluster/namespace) | Hard (assign per resource) |
| **Management overhead** | Low (fewer assignments) | High (more assignments) |
| **Expected latency** | Higher (more DB queries) | Lower (single query) |
| **Use case fit** | Multi-tenant, hierarchical orgs | Simple, flat permissions |

---

## Performance Trade-offs

### Hierarchical Schema

**Advantages:**
- ✅ **50% fewer relationships** to store and manage
- ✅ **Efficient bulk operations**: Assign 1 user at cluster level → affects 10,000 resources
- ✅ **Flexible inheritance**: Different permissions at each level
- ✅ **Better data modeling**: Matches real-world organizational structure

**Disadvantages:**
- ⚠️ **Slower permission checks**: Must traverse 3 levels (resource → namespace → cluster)
- ⚠️ **More database queries**: Each level requires a lookup
- ⚠️ **Higher P99 latency**: Deep traversals can spike latency
- ⚠️ **Complex debugging**: Permission denials harder to trace

**Best for:**
- Multi-tenant systems (cluster = tenant)
- Kubernetes-like resource hierarchies
- Organizations with clear hierarchical structure
- Systems where bulk permission changes are common
- When relationship count is a concern (storage/memory)

### Flat Schema

**Advantages:**
- ✅ **Faster permission checks**: Direct lookup, no traversal
- ✅ **Predictable latency**: Single query per check
- ✅ **Simple debugging**: Easy to see who has access
- ✅ **Lower P99**: No deep traversal spikes

**Disadvantages:**
- ⚠️ **50% more relationships**: More storage and memory
- ⚠️ **Bulk operations expensive**: Must update 10,000 resources individually
- ⚠️ **Management overhead**: More assignments to track
- ⚠️ **Less flexible**: Can't easily group resources

**Best for:**
- Simple permission models
- Systems with minimal hierarchy
- Low-latency requirements (< 5ms)
- When relationship count is not a concern
- When permissions change infrequently

---

## Decision Matrix

| Your Requirement | Recommended Schema |
|------------------|-------------------|
| Multi-tenant architecture | **Hierarchical** |
| Kubernetes-like resources | **Hierarchical** |
| Bulk permission changes common | **Hierarchical** |
| > 100k resources | **Hierarchical** |
| Latency must be < 5ms | **Flat** |
| Simple permission model | **Flat** |
| Permissions change rarely | **Flat** |
| Direct resource access only | **Flat** |

---

## Example Scenarios

### Scenario 1: Add new admin to all resources in cluster0

**Hierarchical:**
```yaml
# 1 relationship
cluster:cluster0#admin@user:new_admin
```
→ Affects 100 namespaces, 10,000 resources instantly

**Flat:**
```yaml
# 10,000 relationships (one per resource)
resource:res0#admin@user:new_admin
resource:res1#admin@user:new_admin
...
resource:res9999#admin@user:new_admin
```
→ Must update each resource individually

### Scenario 2: Check if user can access specific resource

**Hierarchical:**
```
Permission check latency: ~5-15ms (depends on inheritance depth)
- Direct resource permission: ~2ms
- Namespace inherited: ~5ms
- Cluster inherited: ~10ms
```

**Flat:**
```
Permission check latency: ~1-3ms (direct lookup)
- Always direct: ~2ms
```

### Scenario 3: List all resources user can access

**Hierarchical:**
```
If user has cluster-level permission:
- 1 check → returns 10,000 resources
- Fast bulk lookup via inheritance
```

**Flat:**
```
Must check each resource individually:
- 10,000 checks → returns matching resources
- Slower for bulk access
```

---

## Recommended Test Approach

1. **Run the automated test:**
   ```bash
   ./run-performance-test.sh
   ```

2. **Focus on these metrics:**
   - **P95 latency**: What's your tolerance? (5ms? 10ms? 20ms?)
   - **Throughput**: How many checks/sec do you need?
   - **Bulk operation time**: How often do you assign permissions in bulk?

3. **Compare with your use case:**
   - Do you have a natural hierarchy? → Hierarchical
   - Need predictable low latency? → Flat
   - Relationship count matters? → Hierarchical
   - Simple is better? → Flat

4. **Make decision based on:**
   - Your latency requirements
   - Permission management patterns
   - Organizational structure
   - Scalability needs
