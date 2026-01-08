#!/bin/bash

set -e

# Configuration
SPICEDB_PORT=50051
SPICEDB_TOKEN="foobar"
RESULTS_DIR="../performance-results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Test parameters
CLUSTERS=100
NAMESPACES=100
PODS=100

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== SpiceDB Realistic Performance Test ===${NC}"
echo "Testing with your 7 real-world test cases"
echo "Timestamp: $TIMESTAMP"
echo ""
echo "Configuration:"
echo "  Clusters:    $CLUSTERS"
echo "  Namespaces:  $NAMESPACES per cluster"
echo "  Pods:        $PODS per namespace"
echo "  Total:       $(($CLUSTERS * $NAMESPACES * $PODS)) pod resources"
echo ""

# Create results directory
mkdir -p "$RESULTS_DIR"

# Function to start SpiceDB
start_spicedb() {
    echo -e "${YELLOW}Starting SpiceDB...${NC}"

    spicedb serve \
        --grpc-preshared-key "$SPICEDB_TOKEN" \
        --http-enabled=false \
        --datastore-engine=memory \
        --grpc-addr=":$SPICEDB_PORT" &

    SPICEDB_PID=$!
    echo "SpiceDB PID: $SPICEDB_PID"

    # Wait for SpiceDB to be ready
    echo "Waiting for SpiceDB to be ready..."
    for i in {1..30}; do
        if grpcurl -plaintext -d '{}' localhost:$SPICEDB_PORT authzed.api.v1.SchemaService/ReadSchema &>/dev/null 2>&1; then
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

# Cleanup on exit
trap stop_spicedb EXIT

echo -e "${GREEN}Step 1: Generate Test Data${NC}"
echo "=========================================="

echo "Generating realistic test data with 7 test cases..."
go run ../flat/generate-flat-test-data.go \
    --clusters $CLUSTERS \
    --namespaces $NAMESPACES \
    --pods $PODS \
    --test-cases=true \
    --output ../flat/relationships_flat.yaml

echo -e "${GREEN}✓ Test data generated${NC}"
echo ""

echo -e "${GREEN}Step 2: Start SpiceDB${NC}"
echo "=========================================="

start_spicedb
sleep 1

echo ""
echo -e "${GREEN}Step 3: Load Schema${NC}"
echo "=========================================="

echo "Loading realistic schema..."
go run ../common/load-schema.go ../flat/schema_flat.yaml

echo -e "${GREEN}✓ Schema loaded${NC}"
echo ""

echo -e "${GREEN}Step 4: Import Relationships${NC}"
echo "=========================================="

echo "Importing relationships (this may take a while for large datasets)..."

# Import using Go client
go run ../common/load-relationships.go ../flat/relationships_flat.yaml

echo -e "${GREEN}✓ Relationships imported${NC}"
echo ""

echo -e "${GREEN}Step 5: Run Benchmarks${NC}"
echo "=========================================="

RESULT_FILE="$RESULTS_DIR/flat_${TIMESTAMP}.txt"

echo "Running benchmark for 7 test cases..."
go run ../flat/benchmark-flat.go \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --checks 100 \
    --clusters $CLUSTERS \
    --namespaces $NAMESPACES \
    --pods $PODS \
    2>&1 | tee "$RESULT_FILE"

echo -e "${GREEN}✓ Benchmark completed${NC}"
echo ""

echo -e "${GREEN}Step 6: Verify Test Cases${NC}"
echo "=========================================="

echo ""
echo "Verifying each test case with lookup-resources..."

# Test Case 1: user1 can get all resources on all clusters
echo -e "\n${BLUE}Case 1: user1 - All resources on all clusters${NC}"
zed permission lookup-resources resource get user:user1 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure \
    --show-debug-traces=false | head -20

# Test Case 2: user2 can get resources on cluster1, cluster2 only
echo -e "\n${BLUE}Case 2: user2 - Resources on cluster1, cluster2 only${NC}"
zed permission lookup-resources resource get user:user2 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure | grep -E "(cluster1|cluster2)" | head -10

# Test Case 3: user3 can get pods on cluster1 in all namespaces
echo -e "\n${BLUE}Case 3: user3 - Pods on cluster1, all namespaces${NC}"
zed permission lookup-resources resource get user:user3 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure | grep "cluster1" | head -10

# Test Case 4: user4 can get pods on cluster1 in namespace1 only
echo -e "\n${BLUE}Case 4: user4 - Pods on cluster1/namespace1 only${NC}"
zed permission lookup-resources resource get user:user4 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure | head -10

# Test Case 5: user5 can create pods on cluster1 in namespace1
echo -e "\n${BLUE}Case 5: user5 - Create pods on cluster1/namespace1${NC}"
zed permission check resource:cluster/cluster1/namespace/namespace1/core/pods create user:user5 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure

# Test Case 6: user6 can get pods on all clusters
echo -e "\n${BLUE}Case 6: user6 - Pods on all clusters${NC}"
zed permission lookup-resources resource get user:user6 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure | head -10

# Test Case 7: user7 via group1
echo -e "\n${BLUE}Case 7: user7 - Via group1 membership${NC}"
zed permission lookup-resources resource get user:user7 \
    --endpoint localhost:$SPICEDB_PORT \
    --token "$SPICEDB_TOKEN" \
    --insecure | head -10

echo ""
echo -e "${GREEN}=========================================="
echo "Test Complete!"
echo "==========================================${NC}"
echo ""
echo "Results saved to: $RESULT_FILE"
echo ""
echo "Summary:"
cat "$RESULT_FILE" | grep -A 20 "SUMMARY" || echo "Check $RESULT_FILE for full results"

echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Review the results in $RESULT_FILE"
echo "2. Compare with your flat schema results"
echo "3. Analyze which schema works better for your use cases"
echo ""
echo "To compare with flat schema, create flat test data and run similar tests."
