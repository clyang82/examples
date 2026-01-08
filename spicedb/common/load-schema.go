package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type SchemaFile struct {
	Schema string `yaml:"schema"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <schema-file>\n", os.Args[0])
		os.Exit(1)
	}

	fileContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	// Try to parse as YAML first
	var schemaFile SchemaFile
	var schemaContent string

	if err := yaml.Unmarshal(fileContent, &schemaFile); err == nil && schemaFile.Schema != "" {
		// It's a YAML file with schema field
		schemaContent = schemaFile.Schema
	} else {
		// It's a plain schema file
		schemaContent = string(fileContent)
	}

	// Remove leading/trailing whitespace
	schemaContent = strings.TrimSpace(schemaContent)

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
		Schema: schemaContent,
	}

	resp, err := client.WriteSchema(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Schema loaded successfully\n")
	fmt.Printf("  Written at: %s\n", resp.WrittenAt.Token)
}
