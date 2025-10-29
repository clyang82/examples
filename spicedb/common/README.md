# Common Utilities

Shared Go utilities for loading schemas and importing relationships into SpiceDB.

## Files

- `load-schema.go` - Load SpiceDB schema from file
- `load-relationships.go` - Import relationships efficiently

## Usage

### Load Schema

```bash
go run load-schema.go <schema-file>

# Examples
go run load-schema.go ../flat/schema_flat.yaml
go run load-schema.go ../hierarchical/schema_new.yaml
```

**Features:**
- Connects to SpiceDB via gRPC
- Validates schema before loading
- Returns schema write token

### Import Relationships

```bash
go run load-relationships.go <relationships-file>

# Examples
go run load-relationships.go ../flat/relationships_flat.yaml
go run load-relationships.go ../hierarchical/relationships_hierarchical.yaml
```

**Features:**
- Batch imports (1000 relationships per batch)
- Progress tracking
- Error handling
- Performance metrics

**Performance:**
- Flat schema: ~9,000 relationships/sec
- Hierarchical schema: ~11,400 relationships/sec

## Connection Settings

Both utilities connect to SpiceDB using:
- **Endpoint**: `localhost:50051`
- **Token**: `foobar`
- **Transport**: Insecure (for testing)

## Requirements

- SpiceDB running and accessible
- Go 1.19+
- Dependencies from go.mod

## Error Handling

The utilities will exit with status 1 if:
- Cannot connect to SpiceDB
- Schema validation fails
- Relationship import fails
- Invalid file format

## Implementation Details

### Batch Import

Relationships are imported in batches of 1000 to optimize performance:
```go
batchSize := 1000
updates := make([]*v1.RelationshipUpdate, len(relationships))
for i, rel := range relationships {
    updates[i] = &v1.RelationshipUpdate{
        Operation:    v1.RelationshipUpdate_OPERATION_TOUCH,
        Relationship: rel,
    }
}
```

### Relationship Parsing

Supports SpiceDB relationship format:
```
resource:id#relation@subject:subjectId
resource:id#relation@group:groupId#member
```

## Example Output

```
✓ Schema loaded successfully
  Written at: GgYKBENMY0c=

Importing relationships...
  Imported 1000 relationships (10000.0/sec)
  Imported 2000 relationships (10500.0/sec)
  ...
✓ Successfully imported 1000000 relationships in 110s
  Average rate: 9090.9 relationships/sec
```
