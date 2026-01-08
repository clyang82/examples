# Flat Schema Implementation

Simple SpiceDB schema with path-based resource naming and direct permission assignments.

## Files

- `schema_flat.yaml` - Flat schema with path-based resource naming
- `generate-test-data.go` - Generate basic test data
- `generate-flat-test-data.go` - Generate realistic test data with 7 test cases
- `benchmark-flat.go` - Performance benchmark tool

## Schema Design

### Structure

```yaml
definition resource {
  relation admin: user | group#member
  relation viewer: user | group#member
  permission get = admin + viewer
}
```

### Resource Naming Convention

Path-based format: `cluster/{clusterId}/namespace/{nsId}/{apiGroup}/{resourceType}/{name}`

**Examples:**
```
resource:cluster/cluster1/namespace/namespace1/core/pods/pod1
resource:cluster/cluster1/namespace/namespace1/apps/deployments/nginx
```

## Performance Characteristics

- **Avg Latency**: 581µs
- **P95 Latency**: 1.1ms
- **P99 Latency**: 5.8ms
- **Throughput**: 1,780 checks/sec
- **Relationships**: ~1,000,000
- **Import Rate**: ~9,000/sec

## Advantages

✅ Slightly faster permission checks
✅ 50% fewer total relationships
✅ Simple, direct permission model
✅ Lower P99 latency
✅ Easier debugging (direct paths)

## Limitations

❌ Wildcards don't work as pattern matchers
❌ No permission inheritance
❌ Manual per-resource permission assignment
❌ Limited functional coverage (14%)

## Usage

### Generate Test Data

```bash
go run generate-flat-test-data.go \
  --clusters 100 \
  --namespaces 100 \
  --pods 100 \
  --output relationships_flat.yaml
```

### Run Benchmark

```bash
go run benchmark-flat.go \
  --endpoint localhost:50051 \
  --token foobar \
  --checks 1000
```

## When to Use

Choose flat schema if:
- You need minimal latency (< 5ms P99)
- Permissions are managed per-resource
- Simple debugging is priority
- Relationship count matters (storage constrained)
- You implement custom wildcard/pattern matching

## Test Results

See `../docs/FINAL_COMPARISON.md` for complete analysis.
