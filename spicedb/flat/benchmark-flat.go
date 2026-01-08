package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	endpoint    = flag.String("endpoint", "localhost:50051", "SpiceDB gRPC endpoint")
	token       = flag.String("token", "foobar", "SpiceDB token")
	numChecks   = flag.Int("checks", 100, "Number of permission checks per test")
	numClusters = flag.Int("clusters", 100, "Number of clusters in test data")
	numNS       = flag.Int("namespaces", 100, "Number of namespaces per cluster")
	numPods     = flag.Int("pods", 100, "Number of pods per namespace")
)

type TestCase struct {
	Name        string
	Description string
	User        string
	Permission  string
	ResourceIDs []string // Sample resources to check
}

type BenchmarkResult struct {
	Name          string
	TotalDuration time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
	Throughput    float64
	SuccessCount  int
	FailureCount  int
}

func main() {
	flag.Parse()

	client, err := authzed.NewClient(
		*endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(*token),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Printf("=== SpiceDB Realistic Performance Benchmark ===\n")
	fmt.Printf("Endpoint: %s\n", *endpoint)
	fmt.Printf("Test Data: %d clusters × %d namespaces × %d pods = %d total resources\n\n",
		*numClusters, *numNS, *numPods, *numClusters * *numNS * *numPods)

	// Define test cases matching the 7 scenarios
	testCases := []TestCase{
		{
			Name:        "Case 1: Wildcard all clusters",
			Description: "user1 can get all resources on all clusters",
			User:        "user1",
			Permission:  "get",
			ResourceIDs: generateSampleResources(10), // Sample 10 resources
		},
		{
			Name:        "Case 2: Multiple specific clusters",
			Description: "user2 can get resources on cluster1, cluster2 only",
			User:        "user2",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster/cluster1/namespace/namespace0/core/pods/pod0",
				"cluster/cluster2/namespace/namespace0/core/pods/pod0",
				"cluster/cluster3/namespace/namespace0/core/pods/pod0", // Should fail
			},
		},
		{
			Name:        "Case 3: Pods in all namespaces on cluster1",
			Description: "user3 can get pods on cluster1 in all namespaces",
			User:        "user3",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster/cluster1/namespace/namespace0/core/pods/pod0",
				"cluster/cluster1/namespace/namespace50/core/pods/pod50",
				"cluster/cluster2/namespace/namespace0/core/pods/pod0", // Should fail
			},
		},
		{
			Name:        "Case 4: Pods in specific namespace",
			Description: "user4 can get pods on cluster1 in namespace1 only",
			User:        "user4",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster/cluster1/namespace/namespace1/core/pods/pod0",
				"cluster/cluster1/namespace/namespace1/core/pods/pod50",
				"cluster/cluster1/namespace/namespace2/core/pods/pod0", // Should fail
			},
		},
		{
			Name:        "Case 5: Create pods in specific namespace",
			Description: "user5 can create pods on cluster1 in namespace1",
			User:        "user5",
			Permission:  "create",
			ResourceIDs: []string{
				"cluster/cluster1/namespace/namespace1/core/pods",
				"cluster/cluster1/namespace/namespace2/core/pods", // Should fail
			},
		},
		{
			Name:        "Case 6: Specific resource type on all clusters",
			Description: "user6 can get pods on all clusters",
			User:        "user6",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster/cluster0/namespace/namespace0/core/pods/pod0",
				"cluster/cluster50/namespace/namespace50/core/pods/pod50",
				"cluster/cluster99/namespace/namespace99/core/pods/pod99",
			},
		},
		{
			Name:        "Case 7: Group membership",
			Description: "user7 via group1 can get pods on cluster1 in namespace1",
			User:        "user7",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster/cluster1/namespace/namespace1/core/pods/pod0",
				"cluster/cluster1/namespace/namespace1/core/pods/pod50",
				"cluster/cluster1/namespace/namespace2/core/pods/pod0", // Should fail
			},
		},
	}

	// Run benchmarks for each test case
	results := []BenchmarkResult{}
	for _, tc := range testCases {
		fmt.Printf("\n%s\n%s\n", tc.Name, tc.Description)
		fmt.Println(repeatString("-", 60))
		result := runTestCase(ctx, client, tc)
		results = append(results, result)
	}

	// Print summary
	printSummary(results)
}

func runTestCase(ctx context.Context, client *authzed.Client, tc TestCase) BenchmarkResult {
	latencies := []time.Duration{}
	successCount := 0
	failureCount := 0

	start := time.Now()

	// Run multiple checks on sample resources
	for i := 0; i < *numChecks; i++ {
		resourceIdx := i % len(tc.ResourceIDs)
		resourceID := tc.ResourceIDs[resourceIdx]

		checkStart := time.Now()

		req := &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "resource",
				ObjectId:   resourceID,
			},
			Permission: tc.Permission,
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   tc.User,
				},
			},
		}

		resp, err := client.CheckPermission(ctx, req)
		latency := time.Since(checkStart)
		latencies = append(latencies, latency)

		if err != nil {
			failureCount++
			if i == 0 {
				fmt.Printf("  ❌ Error: %v\n", err)
			}
		} else if resp.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			successCount++
		} else {
			// Expected denials for some test resources
			if i < 3 {
				fmt.Printf("  ⚠️  Resource %s: DENIED\n", resourceID)
			}
		}

		if (i+1)%20 == 0 {
			fmt.Printf("  Progress: %d/%d checks\n", i+1, *numChecks)
		}
	}

	totalDuration := time.Since(start)

	result := BenchmarkResult{
		Name:          tc.Name,
		TotalDuration: totalDuration,
		AvgLatency:    avgDuration(latencies),
		P50Latency:    percentile(latencies, 50),
		P95Latency:    percentile(latencies, 95),
		P99Latency:    percentile(latencies, 99),
		Throughput:    float64(*numChecks) / totalDuration.Seconds(),
		SuccessCount:  successCount,
		FailureCount:  failureCount,
	}

	printResult(result)
	return result
}

func generateSampleResources(count int) []string {
	resources := []string{}
	step := (*numClusters * *numNS * *numPods) / count

	for i := 0; i < count; i++ {
		idx := i * step
		clusterIdx := idx / (*numNS * *numPods)
		nsIdx := (idx / *numPods) % *numNS
		podIdx := idx % *numPods

		if clusterIdx >= *numClusters {
			clusterIdx = *numClusters - 1
		}
		if nsIdx >= *numNS {
			nsIdx = *numNS - 1
		}
		if podIdx >= *numPods {
			podIdx = *numPods - 1
		}

		resourceID := fmt.Sprintf("cluster/cluster%d/namespace/namespace%d/core/pods/pod%d",
			clusterIdx, nsIdx, podIdx)
		resources = append(resources, resourceID)
	}

	return resources
}

func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}

func percentile(durations []time.Duration, p int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	idx := (len(sorted) * p) / 100
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func printResult(r BenchmarkResult) {
	fmt.Printf("\n  Results:\n")
	fmt.Printf("    Total Duration: %v\n", r.TotalDuration)
	fmt.Printf("    Avg Latency:    %v\n", r.AvgLatency)
	fmt.Printf("    P50 Latency:    %v\n", r.P50Latency)
	fmt.Printf("    P95 Latency:    %v\n", r.P95Latency)
	fmt.Printf("    P99 Latency:    %v\n", r.P99Latency)
	fmt.Printf("    Throughput:     %.2f checks/sec\n", r.Throughput)
	fmt.Printf("    Success:        %d/%d\n", r.SuccessCount, *numChecks)
	fmt.Printf("    Failures:       %d\n", r.FailureCount)
}

func printSummary(results []BenchmarkResult) {
	fmt.Printf("\n\n")
	fmt.Println(repeatString("=", 80))
	fmt.Println("SUMMARY")
	fmt.Println(repeatString("=", 80))

	fmt.Printf("\n%-40s %10s %10s %10s %12s\n", "Test Case", "Avg", "P95", "P99", "Throughput")
	fmt.Println(repeatString("-", 80))

	for _, r := range results {
		fmt.Printf("%-40s %10v %10v %10v %9.1f/s\n",
			truncate(r.Name, 40),
			r.AvgLatency.Round(time.Microsecond),
			r.P95Latency.Round(time.Microsecond),
			r.P99Latency.Round(time.Microsecond),
			r.Throughput)
	}

	// Calculate overall stats
	var totalAvg, totalP95, totalP99 time.Duration
	var totalThroughput float64

	for _, r := range results {
		totalAvg += r.AvgLatency
		totalP95 += r.P95Latency
		totalP99 += r.P99Latency
		totalThroughput += r.Throughput
	}

	count := len(results)
	fmt.Println(repeatString("-", 80))
	fmt.Printf("%-40s %10v %10v %10v %9.1f/s\n",
		"AVERAGE",
		(totalAvg / time.Duration(count)).Round(time.Microsecond),
		(totalP95 / time.Duration(count)).Round(time.Microsecond),
		(totalP99 / time.Duration(count)).Round(time.Microsecond),
		totalThroughput/float64(count))

	fmt.Println(repeatString("=", 80))
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
