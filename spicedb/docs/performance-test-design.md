# SpiceDB Schema Performance Test Design

## Objective
Compare performance between two schemas:
- **schema_flat.yaml**: Flat structure (resources only)
- **schema_new.yaml**: Hierarchical structure (cluster → namespace → resource)

## Test Configuration
- **100 clusters**
- **100 namespaces** (1 per cluster, or distributed)
- **100 resources** per namespace
- **Total resources**: 100 × 100 × 100 = 1,000,000 relationships

## Key Differences

### schema_flat.yaml (Flat)
- Direct relationships: `resource → user`
- No hierarchy traversal
- Simpler permission evaluation

### schema_new.yaml (Hierarchical)
- Relationships: `cluster → namespace → resource`
- Permission inheritance through `->` operators
- More complex evaluation (traverses hierarchy)

## Data Model Setup

### For schema_new.yaml (Hierarchical)
```
100 clusters (cluster:cluster0 - cluster:cluster99)
├── 100 namespaces total (namespace:ns0 - namespace:ns99)
│   └── Each namespace belongs to 1 cluster
└── 100 resources per namespace = 10,000 total resources
    └── resource:cluster0-ns0-res0, resource:cluster0-ns0-res1, etc.
```

**Relationships to create:**
1. Namespace → Cluster relationships: 100
2. Resource → Namespace relationships: 10,000
3. User permissions (admin/editor/viewer): variable

### For schema_flat.yaml (Flat)
```
10,000 resources (resource:res0 - resource:res9999)
└── Direct user relationships
```

**Relationships to create:**
1. Resource → User relationships: variable

## Test Scenarios

### Scenario 1: Permission Check Performance
Test permission checks at different levels:

**For schema_new.yaml:**
- Check cluster-level permission (affects all resources)
- Check namespace-level permission (affects 100 resources)
- Check resource-level permission (affects 1 resource)

**For schema_flat.yaml:**
- Check resource-level permission only

### Scenario 2: Write Performance
Measure time to create relationships:

1. **Bulk relationship creation**
   - Insert all cluster → namespace relationships
   - Insert all namespace → resource relationships
   - Insert user permission relationships

2. **User assignment at different levels**
   - Assign 10 users as cluster admins
   - Assign 100 users as namespace editors
   - Assign 1000 users as resource viewers

### Scenario 3: Read Performance
Query patterns:

1. **List all resources a user can access**
   - User with cluster-level permissions (should return all 10k resources)
   - User with namespace-level permissions (should return 100 resources)
   - User with resource-level permissions (should return 1 resource)

2. **Check specific permission**
   - `CheckPermission(user:alice, get, resource:cluster0-ns0-res0)`
   - Measure latency across different permission depths

3. **Lookup subjects**
   - Find all users who can `get` a specific resource
   - Compare hierarchy traversal overhead

## Test Data Generation

### Users
```
user:admin1, user:admin2, ..., user:admin10
user:editor1, user:editor2, ..., user:editor100
user:viewer1, user:viewer2, ..., user:viewer1000
```

### Distribution Pattern
```
# Hierarchical (schema_new.yaml)
cluster:cluster0
├── namespace:ns0 (belongs to cluster0)
│   ├── resource:c0-ns0-r0
│   ├── resource:c0-ns0-r1
│   └── ... (100 resources)
├── namespace:ns1 (belongs to cluster0)
└── ... (100 namespaces total across all clusters)

# Flat (schema_flat.yaml)
resource:res0, resource:res1, ..., resource:res9999
```

## Metrics to Measure

1. **Write Operations**
   - Time to insert all relationships
   - Throughput (relationships/second)

2. **Permission Checks**
   - Latency (p50, p95, p99)
   - Throughput (checks/second)

3. **Lookup Operations**
   - Query response time
   - Result set size impact

4. **Resource Usage**
   - Memory consumption
   - Database storage size

## Implementation Steps

### 1. Data Generation Script
Create a script to generate test data for both schemas:
- Generate relationships in SpiceDB format
- Support bulk import

### 2. Load Test Tool
Use SpiceDB's built-in tools or create custom benchmark:
- `spicedb benchmark` command
- Custom Go program using SpiceDB client

### 3. Test Execution
```bash
# For each schema:
1. Start SpiceDB instance
2. Load schema
3. Bulk import relationships
4. Run permission check benchmarks
5. Run lookup benchmarks
6. Collect metrics
```

### 4. Analysis
Compare:
- Permission check latency (hierarchical vs flat)
- Impact of inherited permissions
- Scalability characteristics

## Expected Results

### schema_flat.yaml (Flat) - Expected Advantages
- Faster write operations (fewer relationships)
- Simpler permission checks (no traversal)
- Lower memory overhead

### schema_new.yaml (Hierarchical) - Expected Advantages
- More efficient bulk permission management
- Fewer total relationships for user assignments
- Better organizational model

### Trade-offs
- **Flat**: Fast checks but requires more user assignments per resource
- **Hierarchical**: Slower checks due to traversal but fewer total relationships

## Sample Test Commands

```bash
# Check permission performance
spicedb benchmark check \
  --schema schema_flat.yaml \
  --relationships relationships.yaml \
  --requests 10000 \
  --concurrency 10

# Lookup performance
spicedb benchmark lookup \
  --schema schema_new.yaml \
  --relationships relationships_hierarchical.yaml \
  --requests 1000
```

## Next Steps

1. Create data generation script
2. Generate test data files
3. Run benchmarks
4. Collect and visualize results
5. Make schema decision based on your use case
