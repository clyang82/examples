# Test Automation Scripts

Bash scripts for automating performance tests.

## Files

- `run-flat-test.sh` - Run realistic test suite with 7 test cases
- `run-performance-test.sh` - Run generic flat vs hierarchical comparison
- `compare-results.sh` - Compare test results from two runs

## Usage

### Run Realistic Test

Tests the flat schema with realistic data and 7 test cases:

```bash
./run-flat-test.sh
```

**What it does:**
1. Generates test data (7 test cases + bulk data)
2. Starts SpiceDB
3. Loads schema and relationships
4. Runs benchmarks
5. Verifies each test case
6. Generates results

**Output:**
- `performance-results/flat_TIMESTAMP.txt`
- Test case verification results
- Summary with metrics

### Run Performance Test

Compares flat and hierarchical schemas:

```bash
./run-performance-test.sh
```

**What it does:**
1. Generates data for both schemas
2. Tests hierarchical schema
3. Tests flat schema
4. Generates comparison summary

**Output:**
- `performance-results/hierarchical_TIMESTAMP.txt`
- `performance-results/flat_TIMESTAMP.txt`
- `performance-results/summary_TIMESTAMP.txt`

### Compare Results

Compare two test runs:

```bash
./compare-results.sh <hierarchical-file> <flat-file>

# Example
./compare-results.sh \
  performance-results/hierarchical_20241029_120000.txt \
  performance-results/flat_20241029_120000.txt
```

**Output:**
- Side-by-side comparison
- Key metrics extraction
- Analysis and recommendations

## Configuration

### Test Parameters

Edit scripts to customize:

```bash
# In run-flat-test.sh
CLUSTERS=100
NAMESPACES=100
PODS=100

# In run-performance-test.sh
SPICEDB_PORT=50051
SPICEDB_TOKEN="foobar"
```

### SpiceDB Settings

```bash
SPICEDB_PORT=50051
SPICEDB_TOKEN="foobar"
RESULTS_DIR="./performance-results"
```

## Script Details

### run-flat-test.sh

**Prerequisites:**
- SpiceDB running in Kubernetes
- `kubectl` configured
- Go installed
- Port 50051 available

**Steps:**
1. Generate test data
   ```bash
   go run generate-flat-test-data.go ...
   ```

2. Load schema
   ```bash
   zed schema write ...
   ```

3. Import relationships
   ```bash
   zed relationship import ...
   ```

4. Run benchmarks
   ```bash
   go run benchmark-flat.go ...
   ```

5. Verify test cases
   ```bash
   zed permission lookup-resources ...
   zed permission check ...
   ```

**Duration:** ~10-15 minutes for 1M resources

### run-performance-test.sh

**Prerequisites:**
- SpiceDB installed
- `zed` CLI installed
- Go installed

**Steps:**
1. Generate hierarchical data
2. Start SpiceDB (memory datastore)
3. Load hierarchical schema
4. Import relationships
5. Run benchmark
6. Stop SpiceDB
7. Repeat for flat schema
8. Generate summary

**Duration:** ~20-30 minutes for full comparison

### compare-results.sh

**Prerequisites:**
- Two completed test runs

**Features:**
- Extracts key metrics (avg, P95, P99, throughput)
- Side-by-side comparison
- Recommendations based on results

**Example Output:**
```
=== SpiceDB Schema Performance Comparison ===

Hierarchical Schema Results:
File: performance-results/hierarchical_20241029_120000.txt
----------------------------------------
[Results...]

Flat Schema Results:
File: performance-results/flat_20241029_120000.txt
----------------------------------------
[Results...]

=== Analysis ===
[Comparison and recommendations...]
```

## Environment Variables

### Optional Settings

```bash
# Override SpiceDB endpoint
export SPICEDB_ENDPOINT="localhost:50051"

# Override token
export SPICEDB_TOKEN="your-token"

# Change results directory
export RESULTS_DIR="./custom-results"

# Enable debug mode
export DEBUG=1
```

## Troubleshooting

### SpiceDB Connection Failed

```bash
# Check if SpiceDB is running
kubectl get pods -n spicedb

# Check port-forward
ps aux | grep "port-forward"

# Restart port-forward
pkill -f "port-forward"
kubectl port-forward -n spicedb svc/spicedb 50051:50051 &
```

### Permission Denied

```bash
# Make scripts executable
chmod +x *.sh
```

### Out of Memory

```bash
# Use smaller dataset
# Edit script and reduce:
CLUSTERS=10
NAMESPACES=10
PODS=10
```

### Import Too Slow

```bash
# Use PostgreSQL instead of memory
# Edit k8s-spicedb.yaml:
SPICEDB_DATASTORE_ENGINE: "postgres"
```

## Advanced Usage

### Custom Test Data

```bash
# Edit data generation parameters
go run generate-flat-test-data.go \
  --clusters 50 \
  --namespaces 200 \
  --pods 50 \
  --output custom.yaml
```

### Parallel Testing

```bash
# Run both tests in parallel
./run-flat-test.sh &
PID1=$!

./run-performance-test.sh &
PID2=$!

# Wait for both to complete
wait $PID1 $PID2
```

### Continuous Testing

```bash
# Run tests in a loop
for i in {1..5}; do
  echo "Run $i"
  ./run-flat-test.sh
  sleep 60
done
```

## Output Files

Results are saved in `../performance-results/`:

```
performance-results/
├── flat_20241029_120000.txt
├── hierarchical_20241029_120000.txt
├── flat_20241029_120000.txt
└── summary_20241029_120000.txt
```

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
- name: Run Performance Tests
  run: |
    ./scripts/run-flat-test.sh

- name: Upload Results
  uses: actions/upload-artifact@v2
  with:
    name: performance-results
    path: performance-results/
```

## Customization

### Add New Test Scenario

1. Edit `generate-flat-test-data.go`
2. Add test case relationships
3. Update `benchmark-flat.go`
4. Add test case to benchmark

### Change Metrics

Edit benchmark files to collect additional metrics:
- Query latency
- Error rates
- Resource usage
- Cache hit rates

## Best Practices

1. **Run tests multiple times** for consistency
2. **Use dedicated environment** to avoid interference
3. **Monitor system resources** during tests
4. **Save all results** for comparison
5. **Document test conditions** (load, concurrency, etc.)

## See Also

- [../docs/FLAT_TEST_GUIDE.md](../docs/FLAT_TEST_GUIDE.md) - Detailed test guide
- [../docs/performance-test-design.md](../docs/performance-test-design.md) - Test methodology
- [../docs/FINAL_COMPARISON.md](../docs/FINAL_COMPARISON.md) - Test results
