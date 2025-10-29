package main

import (
	"context"
	"fmt"
	"os"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <schema-file>\n", os.Args[0])
		os.Exit(1)
	}

	schemaContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading schema file: %v\n", err)
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

	req := &v1.WriteSchemaRequest{
		Schema: string(schemaContent),
	}

	resp, err := client.WriteSchema(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Schema loaded successfully\n")
	fmt.Printf("  Written at: %s\n", resp.WrittenAt.Token)
}
