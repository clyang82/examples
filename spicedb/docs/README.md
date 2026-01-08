# Documentation

Comprehensive documentation for SpiceDB performance testing.

## Quick Navigation

| Document | Description |
|----------|-------------|
| [FINAL_COMPARISON.md](FINAL_COMPARISON.md) | ⭐ **Start here** - Complete test results and recommendations |
| [TEST_CASES_ANALYSIS.md](TEST_CASES_ANALYSIS.md) | Detailed analysis of 7 test cases |
| [FLAT_TEST_GUIDE.md](FLAT_TEST_GUIDE.md) | Guide for running realistic tests |
| [SCHEMA_COMPARISON.md](SCHEMA_COMPARISON.md) | Visual comparison of schema designs |
| [performance-test-design.md](performance-test-design.md) | Test methodology and design |
| [PERFORMANCE_TEST_README.md](PERFORMANCE_TEST_README.md) | Original performance test documentation |
| [README_PERFORMANCE_TESTS.md](README_PERFORMANCE_TESTS.md) | Alternative performance test guide |

## Key Findings (TL;DR)

### Winner: Hierarchical Schema ✅

| Metric | Flat | Hierarchical | Winner |
|--------|------|--------------|--------|
| Avg Latency | 581µs | 605µs | Flat (4% faster) |
| P95 Latency | 1.1ms | 881µs | **Hierarchical** (20% faster) |
| Functional Coverage | 14% | 86% | **Hierarchical** |

**Recommendation**: Use hierarchical schema - minimal performance trade-off with vastly better functionality.

## Document Overview

### FINAL_COMPARISON.md
**Complete performance analysis and final recommendation**

Topics covered:
- Executive summary
- Detailed performance metrics
- Schema structure comparison
- Test case functional coverage
- Performance vs. functionality trade-off
- Recommendations and next steps

### TEST_CASES_ANALYSIS.md
**Deep dive into 7 real-world test cases**

Analyzes:
- How each test case works in both schemas
- Relationship count comparison
- Performance differences
- Use case fit analysis
- Migration considerations

### FLAT_TEST_GUIDE.md
**Practical guide for running tests**

Includes:
- Test data structure
- Running benchmarks
- Understanding results
- Customizing tests
- Troubleshooting tips

### SCHEMA_COMPARISON.md
**Visual comparison of schema designs**

Features:
- Data structure diagrams
- Permission evaluation examples
- Trade-off analysis
- Decision matrix
- Example scenarios

### performance-test-design.md
**Test methodology and approach**

Details:
- Test objectives
- Data model setup
- Test scenarios
- Metrics to measure
- Implementation steps

## Test Results Summary

### Performance Metrics

**Flat Schema:**
- Avg: 581µs, P95: 1.1ms, P99: 5.8ms
- Throughput: 1,780 checks/sec
- Relationships: 1,000,119
- Import: 9,000/sec

**Hierarchical Schema:**
- Avg: 605µs, P95: 881µs, P99: 8.5ms
- Throughput: 1,738 checks/sec
- Relationships: 2,020,000
- Import: 11,400/sec

### Functional Coverage

**Flat Schema**: 1/7 test cases work (14%)
**Hierarchical Schema**: 6/7 test cases work (86%)

## Test Cases

1. Global admin on all clusters
2. Multi-cluster admin (cluster1, cluster2)
3. Resource-type scoped (all pods on cluster1)
4. Namespace scoped (pods in specific namespace)
5. Create permission (admin on namespace)
6. Global resource-type (pods on all clusters)
7. Group membership (user via group)

## Reading Guide

### New Users
1. Start with [FINAL_COMPARISON.md](FINAL_COMPARISON.md)
2. Read [TEST_CASES_ANALYSIS.md](TEST_CASES_ANALYSIS.md)
3. Review [SCHEMA_COMPARISON.md](SCHEMA_COMPARISON.md)

### Implementers
1. Read [FLAT_TEST_GUIDE.md](FLAT_TEST_GUIDE.md)
2. Review [performance-test-design.md](performance-test-design.md)
3. Check relevant schema directory (`../flat/` or `../hierarchical/`)

### Researchers
1. Start with [performance-test-design.md](performance-test-design.md)
2. Review [TEST_CASES_ANALYSIS.md](TEST_CASES_ANALYSIS.md)
3. Examine raw test results in `../performance-results/`

## Key Insights

### Why Hierarchical Wins

1. **Permission inheritance works** - Assign at cluster, applies to all resources
2. **All test cases supported** - Flat schema wildcards don't work
3. **Bulk management** - Easier to manage permissions at scale
4. **Only 4% slower** - Minimal performance trade-off
5. **Better P95** - More consistent performance

### Performance Highlights

- Both schemas deliver sub-millisecond average latency
- High throughput: >1,700 checks/second
- Fast imports: 9,000-11,400 relationships/second
- Tested at scale: 1,000,000+ resources

### Important Notes

- SpiceDB wildcards (`_wildcard_`) are literal object IDs, not pattern matchers
- Hierarchical has 2x relationships but better permission management
- P95 latency is better with hierarchical (881µs vs 1.1ms)
- P99 latency is better with flat (5.8ms vs 8.5ms)

## Contributing

To update documentation:
1. Keep summaries concise
2. Include examples and diagrams
3. Link between related documents
4. Update this README when adding new docs

## External Resources

- [SpiceDB Documentation](https://authzed.com/docs)
- [SpiceDB Discord](https://authzed.com/discord)
- [SpiceDB GitHub](https://github.com/authzed/spicedb)
