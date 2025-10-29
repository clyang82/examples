#!/bin/bash

# Simple script to compare performance test results

if [ $# -lt 2 ]; then
    echo "Usage: $0 <hierarchical_result_file> <flat_result_file>"
    echo ""
    echo "Example:"
    echo "  $0 performance-results/hierarchical_20240101_120000.txt performance-results/flat_20240101_120000.txt"
    exit 1
fi

HIERARCHICAL_FILE=$1
FLAT_FILE=$2

if [ ! -f "$HIERARCHICAL_FILE" ]; then
    echo "Error: Hierarchical result file not found: $HIERARCHICAL_FILE"
    exit 1
fi

if [ ! -f "$FLAT_FILE" ]; then
    echo "Error: Flat result file not found: $FLAT_FILE"
    exit 1
fi

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}=== SpiceDB Schema Performance Comparison ===${NC}"
echo ""
echo -e "${BLUE}Hierarchical Schema Results:${NC}"
echo "File: $HIERARCHICAL_FILE"
echo "----------------------------------------"
grep -A 20 "=== Results ===" "$HIERARCHICAL_FILE" || echo "No results found"

echo ""
echo -e "${BLUE}Flat Schema Results:${NC}"
echo "File: $FLAT_FILE"
echo "----------------------------------------"
grep -A 20 "=== Results ===" "$FLAT_FILE" || echo "No results found"

echo ""
echo -e "${YELLOW}=== Analysis ===${NC}"
echo ""

# Extract key metrics for comparison
extract_metric() {
    local file=$1
    local test_name=$2
    local metric=$3

    # Find the test block and extract the metric
    awk -v test="$test_name" -v metric="$metric" '
        /^Test:/ { current_test = $0 }
        current_test ~ test && $0 ~ metric { print $NF }
    ' "$file" | head -1
}

# Compare hierarchical resource-level vs flat direct check
echo "Key Comparison: Resource-level Permission Check"
echo "------------------------------------------------"

HIER_RESOURCE_LAT=$(extract_metric "$HIERARCHICAL_FILE" "Resource-level" "Avg Latency")
FLAT_DIRECT_LAT=$(extract_metric "$FLAT_FILE" "Direct resource" "Avg Latency")

if [ ! -z "$HIER_RESOURCE_LAT" ] && [ ! -z "$FLAT_DIRECT_LAT" ]; then
    echo "Hierarchical (Resource-level):  $HIER_RESOURCE_LAT avg latency"
    echo "Flat (Direct):                  $FLAT_DIRECT_LAT avg latency"
    echo ""
fi

HIER_RESOURCE_P95=$(extract_metric "$HIERARCHICAL_FILE" "Resource-level" "P95 Latency")
FLAT_DIRECT_P95=$(extract_metric "$FLAT_FILE" "Direct resource" "P95 Latency")

if [ ! -z "$HIER_RESOURCE_P95" ] && [ ! -z "$FLAT_DIRECT_P95" ]; then
    echo "P95 Latency:"
    echo "  Hierarchical: $HIER_RESOURCE_P95"
    echo "  Flat:         $FLAT_DIRECT_P95"
    echo ""
fi

HIER_THROUGHPUT=$(extract_metric "$HIERARCHICAL_FILE" "Resource-level" "Throughput")
FLAT_THROUGHPUT=$(extract_metric "$FLAT_FILE" "Direct resource" "Throughput")

if [ ! -z "$HIER_THROUGHPUT" ] && [ ! -z "$FLAT_THROUGHPUT" ]; then
    echo "Throughput:"
    echo "  Hierarchical: $HIER_THROUGHPUT"
    echo "  Flat:         $FLAT_THROUGHPUT"
    echo ""
fi

echo ""
echo -e "${YELLOW}=== Recommendations ===${NC}"
echo ""

echo "Hierarchical Schema Advantages:"
echo "  ✓ Supports permission inheritance (cluster → namespace → resource)"
echo "  ✓ Fewer total relationships to manage"
echo "  ✓ Bulk permission assignment (assign once at cluster/namespace level)"
echo "  ✓ Better organizational model for multi-tenant systems"
echo ""

echo "Flat Schema Advantages:"
echo "  ✓ Simpler permission model"
echo "  ✓ Potentially faster direct permission checks"
echo "  ✓ More predictable latency"
echo "  ✓ Easier debugging"
echo ""

echo "Choose Hierarchical if:"
echo "  • You have a clear organizational hierarchy"
echo "  • You manage permissions at bulk level frequently"
echo "  • Relationship count is a concern"
echo "  • Moderate latency is acceptable"
echo ""

echo "Choose Flat if:"
echo "  • You need minimal latency"
echo "  • Permissions are managed per-resource"
echo "  • Simple permission model is preferred"
echo "  • Relationship count is not a concern"
