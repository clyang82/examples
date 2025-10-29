package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <relationships-file>\n", os.Args[0])
		os.Exit(1)
	}

	client, err := authzed.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("foobar"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Open relationships file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	relationships := []*v1.Relationship{}
	count := 0
	batchSize := 1000
	totalImported := 0
	start := time.Now()

	fmt.Println("Importing relationships...")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines, comments, and the "relationships: |-" header
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "relationships:") {
			continue
		}

		// Parse relationship: resource:id#relation@user:userId
		rel, err := parseRelationship(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid line: %s (error: %v)\n", line, err)
			continue
		}

		relationships = append(relationships, rel)
		count++

		// Write in batches
		if count >= batchSize {
			if err := writeRelationships(ctx, client, relationships); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing batch: %v\n", err)
				os.Exit(1)
			}
			totalImported += count
			fmt.Printf("  Imported %d relationships (%.1f/sec)\n",
				totalImported, float64(totalImported)/time.Since(start).Seconds())
			relationships = []*v1.Relationship{}
			count = 0
		}
	}

	// Write remaining relationships
	if count > 0 {
		if err := writeRelationships(ctx, client, relationships); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing final batch: %v\n", err)
			os.Exit(1)
		}
		totalImported += count
	}

	duration := time.Since(start)
	fmt.Printf("\n✓ Successfully imported %d relationships in %v\n", totalImported, duration)
	fmt.Printf("  Average rate: %.1f relationships/sec\n", float64(totalImported)/duration.Seconds())
}

func parseRelationship(line string) (*v1.Relationship, error) {
	// Format: resource:id#relation@subject:subjectId
	// Example: resource:cluster/cluster1#admin@user:user1

	parts := strings.Split(line, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format: expected '@' separator")
	}

	leftPart := parts[0]
	rightPart := parts[1]

	// Parse left side (resource:id#relation)
	leftSplit := strings.Split(leftPart, "#")
	if len(leftSplit) != 2 {
		return nil, fmt.Errorf("invalid format: expected '#' in resource part")
	}

	resourcePart := leftSplit[0]
	relation := leftSplit[1]

	// Parse resource (type:id)
	resourceSplit := strings.SplitN(resourcePart, ":", 2)
	if len(resourceSplit) != 2 {
		return nil, fmt.Errorf("invalid format: expected ':' in resource")
	}

	resourceType := resourceSplit[0]
	resourceId := resourceSplit[1]

	// Parse right side (subject:id or group:id#member)
	var subjectType, subjectId, subjectRelation string

	if strings.Contains(rightPart, "#") {
		// Group with member relation
		rightSplit := strings.Split(rightPart, "#")
		subjectPart := rightSplit[0]
		subjectRelation = rightSplit[1]

		subjectSplit := strings.SplitN(subjectPart, ":", 2)
		if len(subjectSplit) != 2 {
			return nil, fmt.Errorf("invalid format: expected ':' in subject")
		}
		subjectType = subjectSplit[0]
		subjectId = subjectSplit[1]
	} else {
		// Direct subject
		subjectSplit := strings.SplitN(rightPart, ":", 2)
		if len(subjectSplit) != 2 {
			return nil, fmt.Errorf("invalid format: expected ':' in subject")
		}
		subjectType = subjectSplit[0]
		subjectId = subjectSplit[1]
	}

	rel := &v1.Relationship{
		Resource: &v1.ObjectReference{
			ObjectType: resourceType,
			ObjectId:   resourceId,
		},
		Relation: relation,
		Subject: &v1.SubjectReference{
			Object: &v1.ObjectReference{
				ObjectType: subjectType,
				ObjectId:   subjectId,
			},
		},
	}

	if subjectRelation != "" {
		rel.Subject.OptionalRelation = subjectRelation
	}

	return rel, nil
}

func writeRelationships(ctx context.Context, client *authzed.Client, relationships []*v1.Relationship) error {
	updates := make([]*v1.RelationshipUpdate, len(relationships))
	for i, rel := range relationships {
		updates[i] = &v1.RelationshipUpdate{
			Operation:    v1.RelationshipUpdate_OPERATION_TOUCH,
			Relationship: rel,
		}
	}

	req := &v1.WriteRelationshipsRequest{
		Updates: updates,
	}

	_, err := client.WriteRelationships(ctx, req)
	return err
}
