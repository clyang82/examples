package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	numClusters   = flag.Int("clusters", 100, "Number of clusters")
	numNamespaces = flag.Int("namespaces", 100, "Number of namespaces")
	numResources  = flag.Int("resources", 100, "Number of resources per namespace")
	numAdmins     = flag.Int("admins", 10, "Number of admin users")
	numEditors    = flag.Int("editors", 100, "Number of editor users")
	numViewers    = flag.Int("viewers", 1000, "Number of viewer users")
	schemaType    = flag.String("schema", "hierarchical", "Schema type: 'hierarchical' or 'flat'")
	outputFile    = flag.String("output", "relationships.yaml", "Output file")
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

	if *schemaType == "hierarchical" {
		generateHierarchicalData(file)
	} else {
		generateFlatData(file)
	}

	fmt.Printf("Generated test data in %s\n", *outputFile)
	fmt.Printf("Total relationships: %d\n", countRelationships())
}

func generateHierarchicalData(file *os.File) {
	totalResources := 0

	// Calculate resources per cluster
	resourcesPerCluster := (*numNamespaces * *numResources) / *numClusters
	namespacesPerCluster := *numNamespaces / *numClusters

	fmt.Println("Generating hierarchical data...")
	fmt.Printf("  %d clusters\n", *numClusters)
	fmt.Printf("  %d namespaces (%d per cluster)\n", *numNamespaces, namespacesPerCluster)
	fmt.Printf("  %d resources per namespace\n", *numResources)
	fmt.Printf("  Total resources: %d\n", *numNamespaces * *numResources)

	nsIndex := 0

	// Generate cluster -> namespace -> resource hierarchy
	for c := 0; c < *numClusters; c++ {
		clusterID := fmt.Sprintf("cluster%d", c)

		// Assign cluster-level admins (fewer users, high privilege)
		if c < *numAdmins {
			adminID := fmt.Sprintf("admin%d", c)
			fmt.Fprintf(file, "  cluster:%s#admin@user:%s\n", clusterID, adminID)
		}

		// Create namespaces for this cluster
		for n := 0; n < namespacesPerCluster && nsIndex < *numNamespaces; n++ {
			namespaceID := fmt.Sprintf("ns%d", nsIndex)

			// Link namespace to cluster
			fmt.Fprintf(file, "  namespace:%s#cluster@cluster:%s\n", namespaceID, clusterID)

			// Assign namespace-level editors (moderate users, moderate privilege)
			if nsIndex < *numEditors {
				editorID := fmt.Sprintf("editor%d", nsIndex)
				fmt.Fprintf(file, "  namespace:%s#editor@user:%s\n", namespaceID, editorID)
			}

			// Create resources for this namespace
			for r := 0; r < *numResources; r++ {
				resourceID := fmt.Sprintf("c%d-ns%d-r%d", c, nsIndex, r)

				// Link resource to namespace
				fmt.Fprintf(file, "  resource:%s#namespace@namespace:%s\n", resourceID, namespaceID)

				// Assign resource-level viewers (many users, low privilege)
				viewerIdx := totalResources % *numViewers
				viewerID := fmt.Sprintf("viewer%d", viewerIdx)
				fmt.Fprintf(file, "  resource:%s#viewer@user:%s\n", resourceID, viewerID)

				totalResources++
			}

			nsIndex++
		}
	}

	fmt.Printf("Total resources created: %d\n", totalResources)
}

func generateFlatData(file *os.File) {
	totalResources := *numNamespaces * *numResources

	fmt.Println("Generating flat data...")
	fmt.Printf("  Total resources: %d\n", totalResources)

	// Generate all resources with direct user assignments
	for i := 0; i < totalResources; i++ {
		resourceID := fmt.Sprintf("res%d", i)

		// Assign admins (every resource needs direct admin assignment)
		adminIdx := i % *numAdmins
		adminID := fmt.Sprintf("admin%d", adminIdx)
		fmt.Fprintf(file, "  resource:%s#admin@user:%s\n", resourceID, adminID)

		// Assign editors (every resource needs direct editor assignment)
		editorIdx := i % *numEditors
		editorID := fmt.Sprintf("editor%d", editorIdx)
		fmt.Fprintf(file, "  resource:%s#editor@user:%s\n", resourceID, editorID)

		// Assign viewers (every resource needs direct viewer assignment)
		viewerIdx := i % *numViewers
		viewerID := fmt.Sprintf("viewer%d", viewerIdx)
		fmt.Fprintf(file, "  resource:%s#viewer@user:%s\n", resourceID, viewerID)
	}
}

func countRelationships() int {
	if *schemaType == "hierarchical" {
		// Cluster assignments + namespace->cluster + namespace assignments +
		// resource->namespace + resource assignments
		namespacesPerCluster := *numNamespaces / *numClusters
		totalResources := *numNamespaces * *numResources

		clusterRels := min(*numAdmins, *numClusters)
		namespaceClusterRels := *numNamespaces
		namespaceUserRels := min(*numEditors, *numNamespaces)
		resourceNamespaceRels := totalResources
		resourceUserRels := totalResources // viewers

		return clusterRels + namespaceClusterRels + namespaceUserRels +
			   resourceNamespaceRels + resourceUserRels
	} else {
		// Flat: each resource has 3 direct assignments (admin, editor, viewer)
		totalResources := *numNamespaces * *numResources
		return totalResources * 3
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
