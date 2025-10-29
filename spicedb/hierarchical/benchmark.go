package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	endpoint      = flag.String("endpoint", "localhost:50051", "SpiceDB gRPC endpoint")
	token         = flag.String("token", "foobar", "SpiceDB token")
	schemaType    = flag.String("schema", "hierarchical", "Schema type: 'hierarchical' or 'flat'")
	numChecks     = flag.Int("checks", 1000, "Number of permission checks to perform")
	concurrency   = flag.Int("concurrency", 10, "Number of concurrent workers")
)

type BenchmarkResult struct {
	Name          string
	TotalDuration time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
	Throughput    float64
	Errors        int
}

func main() {
	flag.Parse()

	// Connect to SpiceDB
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

	fmt.Printf("=== SpiceDB Performance Benchmark ===\n")
	fmt.Printf("Schema: %s\n", *schemaType)
	fmt.Printf("Endpoint: %s\n", *endpoint)
	fmt.Printf("Checks: %d\n", *numChecks)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	// Run benchmarks
	results := []BenchmarkResult{}

	if *schemaType == "hierarchical" {
		results = append(results, benchmarkHierarchical(ctx, client)...)
	} else {
		results = append(results, benchmarkFlat(ctx, client)...)
	}

	// Print results
	fmt.Printf("\n=== Results ===\n\n")
	for _, result := range results {
		printResult(result)
	}
}

func benchmarkHierarchical(ctx context.Context, client *authzed.Client) []BenchmarkResult {
	results := []BenchmarkResult{}

	// Test 1: Check permission inherited from cluster level
	fmt.Println("Running: Cluster-level permission check...")
	result1 := runPermissionChecks(ctx, client, "Cluster-level check", func(i int) *v1.CheckPermissionRequest {
		return &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "resource",
				ObjectId:   fmt.Sprintf("c0-ns0-r%d", i%100),
			},
			Permission: "get",
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   "admin0", // Has cluster-level permission
				},
			},
		}
	})
	results = append(results, result1)

	// Test 2: Check permission inherited from namespace level
	fmt.Println("Running: Namespace-level permission check...")
	result2 := runPermissionChecks(ctx, client, "Namespace-level check", func(i int) *v1.CheckPermissionRequest {
		return &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "resource",
				ObjectId:   fmt.Sprintf("c0-ns0-r%d", i%100),
			},
			Permission: "update",
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   "editor0", // Has namespace-level permission
				},
			},
		}
	})
	results = append(results, result2)

	// Test 3: Check permission at resource level
	fmt.Println("Running: Resource-level permission check...")
	result3 := runPermissionChecks(ctx, client, "Resource-level check", func(i int) *v1.CheckPermissionRequest {
		viewerIdx := i % 1000
		resourceIdx := i % 10000
		return &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "resource",
				ObjectId:   fmt.Sprintf("c%d-ns%d-r%d", resourceIdx/10000, (resourceIdx/100)%100, resourceIdx%100),
			},
			Permission: "get",
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   fmt.Sprintf("viewer%d", viewerIdx),
				},
			},
		}
	})
	results = append(results, result3)

	return results
}

func benchmarkFlat(ctx context.Context, client *authzed.Client) []BenchmarkResult {
	results := []BenchmarkResult{}

	// Test: Direct resource permission check
	fmt.Println("Running: Direct resource permission check...")
	result := runPermissionChecks(ctx, client, "Direct resource check", func(i int) *v1.CheckPermissionRequest {
		resourceIdx := i % 10000
		return &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "resource",
				ObjectId:   fmt.Sprintf("res%d", resourceIdx),
			},
			Permission: "get",
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   fmt.Sprintf("viewer%d", i%1000),
				},
			},
		}
	})
	results = append(results, result)

	return results
}

func runPermissionChecks(ctx context.Context, client *authzed.Client, name string, reqGenerator func(int) *v1.CheckPermissionRequest) BenchmarkResult {
	latencies := make([]time.Duration, *numChecks)
	errors := 0

	start := time.Now()

	// Run checks
	for i := 0; i < *numChecks; i++ {
		checkStart := time.Now()

		req := reqGenerator(i)
		_, err := client.CheckPermission(ctx, req)

		latency := time.Since(checkStart)
		latencies[i] = latency

		if err != nil {
			errors++
		}

		// Progress indicator
		if (i+1)%100 == 0 {
			fmt.Printf("  Completed %d/%d checks\n", i+1, *numChecks)
		}
	}

	totalDuration := time.Since(start)

	// Calculate statistics
	return BenchmarkResult{
		Name:          name,
		TotalDuration: totalDuration,
		AvgLatency:    avgDuration(latencies),
		P50Latency:    percentile(latencies, 50),
		P95Latency:    percentile(latencies, 95),
		P99Latency:    percentile(latencies, 99),
		Throughput:    float64(*numChecks) / totalDuration.Seconds(),
		Errors:        errors,
	}
}

func avgDuration(durations []time.Duration) time.Duration {
	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}

func percentile(durations []time.Duration, p int) time.Duration {
	// Simple percentile calculation (should sort for accuracy)
	idx := (len(durations) * p) / 100
	if idx >= len(durations) {
		idx = len(durations) - 1
	}
	return durations[idx]
}

func printResult(r BenchmarkResult) {
	fmt.Printf("Test: %s\n", r.Name)
	fmt.Printf("  Total Duration: %v\n", r.TotalDuration)
	fmt.Printf("  Avg Latency:    %v\n", r.AvgLatency)
	fmt.Printf("  P50 Latency:    %v\n", r.P50Latency)
	fmt.Printf("  P95 Latency:    %v\n", r.P95Latency)
	fmt.Printf("  P99 Latency:    %v\n", r.P99Latency)
	fmt.Printf("  Throughput:     %.2f checks/sec\n", r.Throughput)
	fmt.Printf("  Errors:         %d\n\n", r.Errors)
}
