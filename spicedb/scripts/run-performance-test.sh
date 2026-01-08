#!/bin/bash

set -e

# Configuration
SPICEDB_PORT=50051
SPICEDB_TOKEN="foobar"
RESULTS_DIR="../performance-results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== SpiceDB Performance Test Suite ===${NC}"
echo "Timestamp: $TIMESTAMP"
echo ""

# Create results directory
mkdir -p "$RESULTS_DIR"

# Check if SpiceDB is installed
if ! command -v spicedb &> /dev/null; then
    echo -e "${RED}Error: spicedb command not found${NC}"
    echo "Please install SpiceDB: https://github.com/authzed/spicedb"
    exit 1
fi

# Function to start SpiceDB
start_spicedb() {
    local datastore=$1
    echo -e "${YELLOW}Starting SpiceDB with $datastore datastore...${NC}"

    spicedb serve \
        --grpc-preshared-key "$SPICEDB_TOKEN" \
        --http-enabled=false \
        --datastore-engine=memory \
        --datastore-bootstrap-files="" \
        --grpc-addr=":$SPICEDB_PORT" &

    SPICEDB_PID=$!
    echo "SpiceDB PID: $SPICEDB_PID"

    # Wait for SpiceDB to be ready
    echo "Waiting for SpiceDB to be ready..."
    for i in {1..30}; do
        if grpcurl -plaintext -d '{}' localhost:$SPICEDB_PORT authzed.api.v1.SchemaService/ReadSchema &>/dev/null; then
            echo -e "${GREEN}SpiceDB is ready!${NC}"
            return 0
        fi
        sleep 1
    done

    echo -e "${RED}SpiceDB failed to start${NC}"
    return 1
}

# Function to stop SpiceDB
stop_spicedb() {
    if [ ! -z "$SPICEDB_PID" ]; then
        echo -e "${YELLOW}Stopping SpiceDB (PID: $SPICEDB_PID)...${NC}"
        kill $SPICEDB_PID 2>/dev/null || true
        wait $SPICEDB_PID 2>/dev/null || true
        SPICEDB_PID=""
    fi
}

# Function to load schema
load_schema() {
    local schema_file=$1
    echo -e "${YELLOW}Loading schema from $schema_file...${NC}"

    # Use Go client to load schema
    go run ../common/load-schema.go "$schema_file"

    echo -e "${GREEN}Schema loaded successfully${NC}"
}

# Function to load relationships
load_relationships() {
    local rel_file=$1
    echo -e "${YELLOW}Loading relationships from $rel_file...${NC}"

    # Use Go client to import relationships
    go run ../common/load-relationships.go "$rel_file"

    echo -e "${GREEN}Relationships loaded successfully${NC}"
}

# Function to run benchmark
run_benchmark() {
    local schema_type=$1
    local output_file=$2

    echo -e "${YELLOW}Running benchmark for $schema_type schema...${NC}"

    go run ../hierarchical/benchmark.go \
        --endpoint localhost:$SPICEDB_PORT \
        --token "$SPICEDB_TOKEN" \
        --schema "$schema_type" \
        --checks 1000 \
        --concurrency 10 \
        2>&1 | tee "$output_file"

    echo -e "${GREEN}Benchmark completed${NC}"
}

# Cleanup on exit
trap stop_spicedb EXIT

echo -e "${GREEN}Step 1: Generate test data${NC}"
echo "----------------------------------------"

# Generate data for hierarchical schema
echo "Generating hierarchical test data..."
go run ../hierarchical/generate-hierarchical-data.go \
    --clusters 100 \
    --namespaces 100 \
    --pods 100 \
    --output ../hierarchical/relationships_hierarchical.yaml

# Generate data for flat schema
echo "Generating flat test data..."
go run ../flat/generate-test-data.go \
    --schema flat \
    --clusters 100 \
    --namespaces 100 \
    --resources 100 \
    --admins 10 \
    --editors 100 \
    --viewers 1000 \
    --output ../flat/relationships_flat.yaml

echo ""
echo -e "${GREEN}Step 2: Test Hierarchical Schema (schema_new.yaml)${NC}"
echo "----------------------------------------"

# Start SpiceDB
start_spicedb "memory"

# Load hierarchical schema
load_schema "../hierarchical/schema_new.yaml"

# Load relationships
load_relationships "../hierarchical/relationships_hierarchical.yaml"

# Run benchmark
HIERARCHICAL_RESULT="$RESULTS_DIR/hierarchical_${TIMESTAMP}.txt"
run_benchmark "hierarchical" "$HIERARCHICAL_RESULT"

# Stop SpiceDB
stop_spicedb
sleep 2

echo ""
echo -e "${GREEN}Step 3: Test Flat Schema (schema_flat.yaml)${NC}"
echo "----------------------------------------"

# Start fresh SpiceDB instance
start_spicedb "memory"

# Load flat schema
load_schema "../flat/schema_flat.yaml"

# Load relationships
load_relationships "../flat/relationships_flat.yaml"

# Run benchmark
FLAT_RESULT="$RESULTS_DIR/flat_${TIMESTAMP}.txt"
run_benchmark "flat" "$FLAT_RESULT"

# Stop SpiceDB
stop_spicedb

echo ""
echo -e "${GREEN}=== Test Complete ===${NC}"
echo "----------------------------------------"
echo "Results saved to:"
echo "  Hierarchical: $HIERARCHICAL_RESULT"
echo "  Flat:         $FLAT_RESULT"
echo ""
echo "To compare results, check the files in $RESULTS_DIR"

# Generate summary
SUMMARY_FILE="$RESULTS_DIR/summary_${TIMESTAMP}.txt"
{
    echo "=== Performance Test Summary ==="
    echo "Date: $(date)"
    echo ""
    echo "Configuration:"
    echo "  - 100 clusters"
    echo "  - 100 namespaces"
    echo "  - 100 resources per namespace"
    echo "  - Total: 10,000 resources"
    echo ""
    echo "=== Hierarchical Schema Results ==="
    grep -A 8 "=== Results ===" "$HIERARCHICAL_RESULT" | tail -n +2
    echo ""
    echo "=== Flat Schema Results ==="
    grep -A 8 "=== Results ===" "$FLAT_RESULT" | tail -n +2
} > "$SUMMARY_FILE"

echo "Summary: $SUMMARY_FILE"
cat "$SUMMARY_FILE"
