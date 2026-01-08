package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	numClusters      = flag.Int("clusters", 100, "Number of clusters")
	numNamespaces    = flag.Int("namespaces", 100, "Number of namespaces per cluster")
	numPods          = flag.Int("pods", 100, "Number of pods per namespace")
	outputFile       = flag.String("output", "relationships_flat.yaml", "Output file")
	includeTestCases = flag.Bool("test-cases", true, "Include the 7 specific test cases")
)

func main() {
	flag.Parse()

	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(file, "relationships: |-\n")

	if *includeTestCases {
		generateTestCases(file)
	}

	generateBulkData(file)

	fmt.Printf("Generated realistic test data in %s\n", *outputFile)
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Clusters: %d\n", *numClusters)
	fmt.Printf("  Namespaces per cluster: %d\n", *numNamespaces)
	fmt.Printf("  Pods per namespace: %d\n", *numPods)
	fmt.Printf("  Total resources: %d\n", *numClusters * *numNamespaces * *numPods)
}

func generateTestCases(file *os.File) {
	fmt.Println("\nGenerating specific test cases...")

	// Case 1: user1 can get all resources on all clusters
	// Uses wildcard at cluster level
	fmt.Fprintf(file, "  // Case 1: user1 can get all resources on all clusters\n")
	fmt.Fprintf(file, "  resource:cluster/_wildcard_#admin@user:user1\n\n")

	// Case 2: user2 can only get all resources on cluster1, cluster2 but not cluster3
	fmt.Fprintf(file, "  // Case 2: user2 can get all resources on cluster1, cluster2 but not cluster3\n")
	fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/_wildcard_#admin@user:user2\n")
	fmt.Fprintf(file, "  resource:cluster/cluster2/namespace/_wildcard_#admin@user:user2\n\n")

	// Case 3: user3 can get pods on cluster1 in all the namespaces
	fmt.Fprintf(file, "  // Case 3: user3 can get pods on cluster1 in all namespaces\n")
	fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/_wildcard_/core/pods/_wildcard_#viewer@user:user3\n\n")

	// Case 4: User4 can get pods on cluster1 in namespace1 but not other namespaces
	fmt.Fprintf(file, "  // Case 4: user4 can get pods on cluster1 in namespace1 only\n")
	fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/namespace1/core/pods/_wildcard_#viewer@user:user4\n\n")

	// Case 5: User5 can create pod on cluster1 in namespace1
	fmt.Fprintf(file, "  // Case 5: user5 can create pods on cluster1 in namespace1\n")
	fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/namespace1/core/pods#admin@user:user5\n\n")

	// Case 6: User6 can get pods on all clusters but not other resources
	fmt.Fprintf(file, "  // Case 6: user6 can get pods on all clusters\n")
	fmt.Fprintf(file, "  resource:cluster/_wildcard_/namespace/_wildcard_/core/pods/_wildcard_#viewer@user:user6\n\n")

	// Case 7: group1 can get pod on cluster1 in namespace1, user7 is in group1
	fmt.Fprintf(file, "  // Case 7: group1 has access, user7 is member of group1\n")
	fmt.Fprintf(file, "  group:group1#member@user:user7\n")
	fmt.Fprintf(file, "  resource:cluster/cluster1/namespace/namespace1/core/pods/_wildcard_#viewer@group:group1#member\n\n")

	fmt.Println("  ✓ Generated 7 test cases")
}

func generateBulkData(file *os.File) {
	fmt.Println("\nGenerating bulk test data...")

	totalResources := 0

	// Generate realistic cluster/namespace/pod structure
	for c := 0; c < *numClusters; c++ {
		clusterID := fmt.Sprintf("cluster%d", c)

		// Create cluster-level wildcard resource for some users
		if c < 5 {
			// Some users have cluster-level admin on specific clusters
			fmt.Fprintf(file, "  resource:cluster/%s/namespace/_wildcard_#admin@user:cluster_admin%d\n", clusterID, c)
		}

		for ns := 0; ns < *numNamespaces; ns++ {
			namespaceID := fmt.Sprintf("namespace%d", ns)

			// Create namespace-level wildcard for some users
			if ns < 10 {
				fmt.Fprintf(file, "  resource:cluster/%s/namespace/%s/_wildcard_#editor@user:namespace_editor%d\n",
					clusterID, namespaceID, ns)
			}

			// Generate individual pods
			for p := 0; p < *numPods; p++ {
				podID := fmt.Sprintf("pod%d", p)
				resourcePath := fmt.Sprintf("cluster/%s/namespace/%s/core/pods/%s", clusterID, namespaceID, podID)

				// Assign viewers to specific pods (distributed)
				viewerIdx := totalResources % 1000
				fmt.Fprintf(file, "  resource:%s#viewer@user:pod_viewer%d\n", resourcePath, viewerIdx)

				totalResources++

				// Progress indicator
				if totalResources%10000 == 0 {
					fmt.Printf("  Generated %d resources...\n", totalResources)
				}
			}
		}
	}

	fmt.Printf("  ✓ Generated %d total pod resources\n", totalResources)
}
