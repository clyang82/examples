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
	endpoint  = flag.String("endpoint", "localhost:50051", "SpiceDB gRPC endpoint")
	token     = flag.String("token", "foobar", "SpiceDB token")
	numChecks = flag.Int("checks", 100, "Number of permission checks per test")
)

type TestCase struct {
	Name        string
	Description string
	User        string
	Permission  string
	ResourceIDs []string
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

	fmt.Printf("=== SpiceDB Hierarchical Performance Benchmark ===\n")
	fmt.Printf("Endpoint: %s\n", *endpoint)
	fmt.Printf("Testing hierarchical schema with inheritance\n\n")

	testCases := []TestCase{
		{
			Name:        "Case 1: Admin on all clusters",
			Description: "user1 has admin on all 100 clusters (inherits to all resources)",
			User:        "user1",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster0/namespace0/pods/pod0",
				"cluster0/namespace50/pods/pod50",
				"cluster50/namespace0/pods/pod0",
				"cluster99/namespace99/pods/pod99",
			},
		},
		{
			Name:        "Case 2: Admin on specific clusters",
			Description: "user2 has admin on cluster1, cluster2 only",
			User:        "user2",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster1/namespace0/pods/pod0",
				"cluster2/namespace0/pods/pod0",
				"cluster3/namespace0/pods/pod0", // Should DENY
			},
		},
		{
			Name:        "Case 3: Viewer on cluster1",
			Description: "user3 has viewer on cluster1 (all namespaces/pods)",
			User:        "user3",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster1/namespace0/pods/pod0",
				"cluster1/namespace50/pods/pod50",
				"cluster2/namespace0/pods/pod0", // Should DENY
			},
		},
		{
			Name:        "Case 4: Viewer on specific namespace",
			Description: "user4 has viewer on cluster1/namespace1 only",
			User:        "user4",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster1/namespace1/pods/pod0",
				"cluster1/namespace1/pods/pod50",
				"cluster1/namespace2/pods/pod0", // Should DENY
			},
		},
		{
			Name:        "Case 5: Admin on namespace (create permission)",
			Description: "user5 has admin on cluster1/namespace1",
			User:        "user5",
			Permission:  "create",
			ResourceIDs: []string{
				"cluster1/namespace1/pods/pod0",
				"cluster1/namespace2/pods/pod0", // Should DENY
			},
		},
		{
			Name:        "Case 6: Viewer on all clusters",
			Description: "user6 has viewer on all 100 clusters",
			User:        "user6",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster0/namespace0/pods/pod0",
				"cluster50/namespace50/pods/pod50",
				"cluster99/namespace99/pods/pod99",
			},
		},
		{
			Name:        "Case 7: Group membership",
			Description: "user7 via group1 has viewer on cluster1/namespace1",
			User:        "user7",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster1/namespace1/pods/pod0",
				"cluster1/namespace1/pods/pod50",
				"cluster1/namespace2/pods/pod0", // Should DENY
			},
		},
		{
			Name:        "Case 8: Cluster-scoped resources (nodes)",
			Description: "user1 has admin on all clusters (can access cluster-scoped nodes)",
			User:        "user1",
			Permission:  "get",
			ResourceIDs: []string{
				"cluster0/nodes/node0",
				"cluster1/nodes/node5",
				"cluster50/persistentvolumes/pv0",
			},
		},
	}

	results := []BenchmarkResult{}
	for _, tc := range testCases {
		fmt.Printf("\n%s\n%s\n", tc.Name, tc.Description)
		fmt.Println(repeatString("-", 60))
		result := runTestCase(ctx, client, tc)
		results = append(results, result)
	}

	printSummary(results)
}

func runTestCase(ctx context.Context, client *authzed.Client, tc TestCase) BenchmarkResult {
	latencies := []time.Duration{}
	successCount := 0
	failureCount := 0

	start := time.Now()

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
