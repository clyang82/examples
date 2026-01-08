# SpiceDB Performance Testing Suite

Complete performance testing framework for comparing SpiceDB schema designs with 1,000,000+ resources.

## 🎯 Project Overview

This repository contains a comprehensive performance testing suite for SpiceDB that compares two schema approaches:
- **Flat Schema**: Simple, path-based resource naming
- **Hierarchical Schema**: Multi-level inheritance (cluster → namespace → resource)

**Test Scale**: 100 clusters × 100 namespaces × 100 pods = **1,000,000 resources**

## 📊 Test Results Summary

| Metric | Flat Schema | Hierarchical Schema | Winner |
|--------|-------------|---------------------|--------|
| Avg Latency | 581µs | 605µs | Flat (4% faster) |
| P95 Latency | 1.1ms | 881µs | **Hierarchical** (20% faster) |
| P99 Latency | 5.8ms | 8.5ms | Flat (32% faster) |
| Throughput | 1,780/s | 1,738/s | Flat (2% higher) |
| Functional Coverage | 14% | **86%** | **Hierarchical** |
| **Recommendation** | - | ✅ **Use This** | **Hierarchical** |

**Winner: Hierarchical Schema** - Better functionality with minimal performance trade-off

## 📁 Directory Structure

```
spicedb/
├── flat/                    # Flat schema implementation
│   ├── schema_flat.yaml # Flat schema with path-based naming
│   ├── generate-*.go        # Data generators
│   └── benchmark-*.go       # Benchmarks
│
├── hierarchical/            # Hierarchical schema implementation
│   ├── schema_new.yaml      # Schema with inheritance
│   ├── generate-*.go        # Data generators
│   └── benchmark-*.go       # Benchmarks
│
├── common/                  # Shared utilities
│   ├── load-schema.go       # Schema loader
│   └── load-relationships.go # Relationship importer
│
├── manifests/               # K8s deployment manifests
│   ├── k8s-postgres.yaml    # PostgreSQL deployment
│   └── k8s-spicedb.yaml     # SpiceDB deployment
│
├── scripts/                 # Test automation scripts
│   ├── run-flat-test.sh
│   ├── run-performance-test.sh
│   └── compare-results.sh
│
├── docs/                    # Documentation
│   ├── FINAL_COMPARISON.md  # Complete test results
│   ├── TEST_CASES_ANALYSIS.md
│   └── ...
│
└── performance-results/     # Test output files
```

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (KinD recommended)
- Go 1.19+
- kubectl

### Run Tests

```bash
# 1. Deploy SpiceDB on Kubernetes
kubectl apply -f manifests/k8s-postgres.yaml
kubectl wait --for=condition=ready pod -l app=postgres -n spicedb --timeout=120s
kubectl apply -f manifests/k8s-spicedb.yaml

# 2. Port-forward to SpiceDB
kubectl port-forward -n spicedb svc/spicedb 50051:50051 &

# 3. Run hierarchical schema test
cd hierarchical
go run generate-hierarchical-data.go --clusters 100 --namespaces 100 --pods 100
go run ../common/load-schema.go schema_new.yaml
go run ../common/load-relationships.go relationships_hierarchical.yaml
go run benchmark-hierarchical.go

# 4. Run flat schema test
cd ../flat
go run generate-flat-test-data.go --clusters 100 --namespaces 100 --pods 100
go run ../common/load-schema.go schema_flat.yaml
go run ../common/load-relationships.go relationships_flat.yaml
go run benchmark-flat.go
```

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [FINAL_COMPARISON.md](docs/FINAL_COMPARISON.md) | Complete performance comparison and analysis |
| [TEST_CASES_ANALYSIS.md](docs/TEST_CASES_ANALYSIS.md) | Detailed analysis of 7 test cases |
| [FLAT_TEST_GUIDE.md](docs/FLAT_TEST_GUIDE.md) | Guide for realistic testing |
| [SCHEMA_COMPARISON.md](docs/SCHEMA_COMPARISON.md) | Schema design comparison |

## 🧪 Test Cases

The test suite evaluates 7 real-world scenarios:

1. **Global Admin**: User with access to all clusters
2. **Multi-Cluster Admin**: User with access to specific clusters
3. **Resource-Type Scoped**: User with access to all pods on cluster1
4. **Namespace Scoped**: User with access to specific namespace
5. **Create Permission**: User with create permission on specific namespace
6. **Global Resource-Type**: User with access to pods on all clusters
7. **Group Membership**: User accessing via group membership

## 🏆 Key Findings

### Why Hierarchical Schema Wins

✅ **All 7 test cases work** (vs. 1 for flat)
✅ **Permission inheritance** (cluster → namespace → resource)
✅ **Bulk permission management** via hierarchy
✅ **Only 4% slower** on average (605µs vs 581µs)
✅ **Better P95 latency** (881µs vs 1.1ms)

### Performance Highlights

- **Sub-millisecond latency**: Both schemas deliver <1ms average
- **High throughput**: >1,700 checks/second
- **Fast imports**: 9,000-11,400 relationships/second
- **Scalable**: Tested with 1M+ resources

## 🔧 Components

### Flat Schema (`flat/`)
- Simple resource definition
- Path-based naming: `cluster/cluster1/namespace/ns1/core/pods/pod1`
- Direct relationship assignments
- No permission inheritance

### Hierarchical Schema (`hierarchical/`)
- Multi-level definitions: cluster, namespace, resource
- Permission inheritance via `->` operator
- Resource naming: `cluster1-namespace1-pod1`
- Bulk permission management

### Common Utilities (`common/`)
- `load-schema.go`: Load SpiceDB schemas
- `load-relationships.go`: Import relationships efficiently (1000/batch)

### Kubernetes Deployment (`manifests/`)
- PostgreSQL backend for SpiceDB
- SpiceDB server with migrations
- Service definitions

## 📊 Performance Results

### Import Performance

| Schema | Relationships | Import Time | Rate |
|--------|---------------|-------------|------|
| Flat | 1,000,119 | 110s | 9,000/s |
| Hierarchical | 2,020,000 | 177s | 11,400/s |

### Permission Check Performance

| Schema | Avg | P50 | P95 | P99 | Throughput |
|--------|-----|-----|-----|-----|------------|
| Flat | 581µs | 450µs | 1.1ms | 5.8ms | 1,780/s |
| Hierarchical | 605µs | 460µs | 881µs | 8.5ms | 1,738/s |

## 🛠️ Usage Examples

### Generate Test Data

```bash
# Hierarchical schema
cd hierarchical
go run generate-hierarchical-data.go \
  --clusters 100 \
  --namespaces 100 \
  --pods 100 \
  --output relationships.yaml

# Flat schema
cd flat
go run generate-flat-test-data.go \
  --clusters 100 \
  --namespaces 100 \
  --pods 100 \
  --output relationships.yaml
```

### Run Benchmarks

```bash
# Hierarchical
cd hierarchical
go run benchmark-hierarchical.go \
  --endpoint localhost:50051 \
  --token foobar \
  --checks 1000

# Flat
cd flat
go run benchmark-flat.go \
  --endpoint localhost:50051 \
  --token foobar \
  --checks 1000
```

## 📝 Notes

- **Wildcard Limitation**: Flat schema wildcards (`_wildcard_`) are literal object IDs in SpiceDB, not pattern matchers
- **Relationship Count**: Hierarchical has 2x more relationships due to hierarchy links
- **Use Case Fit**: Hierarchical better matches Kubernetes-like organizational models

## 🤝 Contributing

This is a testing framework. To extend:

1. Add new test cases in benchmark files
2. Create new schema variations in respective directories
3. Update documentation in `docs/`

## 📄 License

See parent repository license.

## 🔗 Resources

- [SpiceDB Documentation](https://authzed.com/docs)
- [Complete Test Results](docs/FINAL_COMPARISON.md)
- [Test Case Analysis](docs/TEST_CASES_ANALYSIS.md)
