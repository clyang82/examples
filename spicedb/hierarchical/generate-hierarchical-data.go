package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	numClusters          = flag.Int("clusters", 100, "Number of clusters")
	numNamespaces        = flag.Int("namespaces", 100, "Number of namespaces per cluster")
	numPods              = flag.Int("pods", 100, "Number of pods per namespace")
	numClusterResources  = flag.Int("cluster-resources", 10, "Number of cluster-scoped resources per cluster (nodes, pvs)")
	outputFile           = flag.String("output", "relationships_hierarchical.yaml", "Output file")
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

	fmt.Println("\nGenerating hierarchical test data...")
	fmt.Printf("  %d clusters\n", *numClusters)
	fmt.Printf("  %d namespaces per cluster\n", *numNamespaces)
	fmt.Printf("  %d pods per namespace\n", *numPods)
	fmt.Printf("  %d cluster-scoped resources per cluster\n", *numClusterResources)

	// Generate 7 test cases
	generateTestCases(file)

	// Generate bulk hierarchical data
	generateHierarchy(file)

	fmt.Printf("\nGenerated hierarchical test data in %s\n", *outputFile)
}

func generateTestCases(file *os.File) {
	fmt.Println("\nGenerating 7 test cases...")

	// Case 1: user1 can get all resources on all clusters
	// Assign user1 as admin on all clusters
	fmt.Fprintf(file, "  // Case 1: user1 - admin on all clusters\n")
	for c := 0; c < *numClusters; c++ {
		clusterID := fmt.Sprintf("cluster%d", c)
		fmt.Fprintf(file, "  cluster:%s#admin@user:user1\n", clusterID)
	}
	fmt.Fprintf(file, "\n")

	// Case 2: user2 can get resources on cluster1, cluster2 only
	fmt.Fprintf(file, "  // Case 2: user2 - admin on cluster1, cluster2\n")
	fmt.Fprintf(file, "  cluster:cluster1#admin@user:user2\n")
	fmt.Fprintf(file, "  cluster:cluster2#admin@user:user2\n\n")

	// Case 3: user3 can get pods on cluster1 in all namespaces
	// Assign user3 as viewer on cluster1
	fmt.Fprintf(file, "  // Case 3: user3 - viewer on cluster1 (all namespaces)\n")
	fmt.Fprintf(file, "  cluster:cluster1#viewer@user:user3\n\n")

	// Case 4: user4 can get pods on cluster1 in namespace1 only
	fmt.Fprintf(file, "  // Case 4: user4 - viewer on namespace1\n")
	fmt.Fprintf(file, "  namespace:cluster1/namespace1#viewer@user:user4\n\n")

	// Case 5: user5 can create pods on cluster1 in namespace1
	fmt.Fprintf(file, "  // Case 5: user5 - admin on namespace1\n")
	fmt.Fprintf(file, "  namespace:cluster1/namespace1#admin@user:user5\n\n")

	// Case 6: user6 can get pods on all clusters
	// Assign user6 as viewer on all clusters
	fmt.Fprintf(file, "  // Case 6: user6 - viewer on all clusters (all pods)\n")
	for c := 0; c < *numClusters; c++ {
		clusterID := fmt.Sprintf("cluster%d", c)
		fmt.Fprintf(file, "  cluster:%s#viewer@user:user6\n", clusterID)
	}
	fmt.Fprintf(file, "\n")

	// Case 7: user7 via group1
	fmt.Fprintf(file, "  // Case 7: user7 - member of group1, group1 viewer on namespace1\n")
	fmt.Fprintf(file, "  group:group1#user@user:user7\n")
	fmt.Fprintf(file, "  namespace:cluster1/namespace1#viewer@group:group1#member\n\n")

	fmt.Println("  ✓ Generated 7 test cases")
}

func generateHierarchy(file *os.File) {
	fmt.Println("\nGenerating hierarchical structure...")

	totalNamespaces := 0
	totalResources := 0
	totalClusterResources := 0

	for c := 0; c < *numClusters; c++ {
		clusterID := fmt.Sprintf("cluster%d", c)

		// Assign some cluster-level admins
		if c < 5 {
			adminID := fmt.Sprintf("cluster_admin%d", c)
			fmt.Fprintf(file, "  cluster:%s#admin@user:%s\n", clusterID, adminID)
		}

		// Generate cluster-scoped resources (nodes, persistentvolumes)
		for r := 0; r < *numClusterResources; r++ {
			// Generate nodes
			nodeID := fmt.Sprintf("%s/nodes/node%d", clusterID, r)
			fmt.Fprintf(file, "  resource:%s#cluster@cluster:%s\n", nodeID, clusterID)

			// Assign some direct permissions on nodes
			if r%5 == 0 {
				viewerID := fmt.Sprintf("node_viewer%d", r)
				fmt.Fprintf(file, "  resource:%s#viewer@user:%s\n", nodeID, viewerID)
			}

			// Generate persistent volumes
			pvID := fmt.Sprintf("%s/persistentvolumes/pv%d", clusterID, r)
			fmt.Fprintf(file, "  resource:%s#cluster@cluster:%s\n", pvID, clusterID)

			totalClusterResources += 2 // node + pv
		}

		// Generate namespaces for this cluster
		for ns := 0; ns < *numNamespaces; ns++ {
			namespaceID := fmt.Sprintf("%s/namespace%d", clusterID, ns)

			// Link namespace to cluster
			fmt.Fprintf(file, "  namespace:%s#cluster@cluster:%s\n", namespaceID, clusterID)

			// Assign namespace-level editors (first 10 namespaces per cluster)
			if ns < 10 {
				editorID := fmt.Sprintf("namespace_editor%d", ns)
				fmt.Fprintf(file, "  namespace:%s#editor@user:%s\n", namespaceID, editorID)
			}

			totalNamespaces++

			// Generate resources (pods) for this namespace
			for p := 0; p < *numPods; p++ {
				resourceID := fmt.Sprintf("%s/pods/pod%d", namespaceID, p)

				// Link resource to namespace
				fmt.Fprintf(file, "  resource:%s#namespace@namespace:%s\n", resourceID, namespaceID)

				// Assign resource-level viewers (distributed)
				viewerIdx := totalResources % 1000
				viewerID := fmt.Sprintf("pod_viewer%d", viewerIdx)
				fmt.Fprintf(file, "  resource:%s#viewer@user:%s\n", resourceID, viewerID)

				totalResources++

				if totalResources%10000 == 0 {
					fmt.Printf("  Generated %d resources...\n", totalResources)
				}
			}
		}
	}

	fmt.Printf("  ✓ Generated %d cluster-scoped resources (nodes + pvs)\n", totalClusterResources)
	fmt.Printf("  ✓ Generated %d namespaces\n", totalNamespaces)
	fmt.Printf("  ✓ Generated %d namespace-scoped resources (pods)\n", totalResources)
}
