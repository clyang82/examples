# Test Queries for ClusterRole RBAC Schema

This file contains example queries to test the SpiceDB schema for Kubernetes ClusterRole permissions.

## Prerequisites

1. SpiceDB server running
2. Schema loaded from `schema.yaml`
3. Relationships loaded from `relationships.yaml`

## Setup

```bash
# Start SpiceDB (if not already running)
spicedb serve --grpc-preshared-key "somerandomkeyhere" &

# Load schema
spicedb schema write --schema schema.yaml --endpoint localhost:50051 --token "somerandomkeyhere"

# Load relationships
cat relationships.yaml | spicedb relationship import --endpoint localhost:50051 --token "somerandomkeyhere"
```

## Test Cases

### 1. Basic Resource Type Permissions

#### Test: Alice can get virtualmachines
```bash
spicedb relationship check \
  resourcetype:virtualmachines get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can delete virtualmachines
```bash
spicedb relationship check \
  resourcetype:virtualmachines delete user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can create virtualmachineinstances
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstances create user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update virtualmachineinstancepresets
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstancepresets update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can patch virtualmachineinstancereplicasets
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstancereplicasets patch user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can list virtualmachineinstancemigrations
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstancemigrations list user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can watch virtualmachines
```bash
spicedb relationship check \
  resourcetype:virtualmachines watch user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can deletecollection virtualmachines
```bash
spicedb relationship check \
  resourcetype:virtualmachines deletecollection user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 2. Snapshot Resources

#### Test: Alice can create virtualmachinesnapshots
```bash
spicedb relationship check \
  resourcetype:virtualmachinesnapshots create user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can delete virtualmachinesnapshotcontents
```bash
spicedb relationship check \
  resourcetype:virtualmachinesnapshotcontents delete user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update virtualmachinerestores
```bash
spicedb relationship check \
  resourcetype:virtualmachinerestores update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 3. Export, Clone, and Instance Type Resources

#### Test: Alice can create virtualmachineexports
```bash
spicedb relationship check \
  resourcetype:virtualmachineexports create user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can delete virtualmachineclones
```bash
spicedb relationship check \
  resourcetype:virtualmachineclones delete user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update virtualmachineinstancetypes
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstancetypes update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can patch virtualmachineclusterinstancetypes
```bash
spicedb relationship check \
  resourcetype:virtualmachineclusterinstancetypes patch user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can list virtualmachinepreferences
```bash
spicedb relationship check \
  resourcetype:virtualmachinepreferences list user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can watch virtualmachineclusterpreferences
```bash
spicedb relationship check \
  resourcetype:virtualmachineclusterpreferences watch user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 4. Pool Resources

#### Test: Alice can create virtualmachinepools
```bash
spicedb relationship check \
  resourcetype:virtualmachinepools create user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 5. Migration Policies (Read-only)

#### Test: Alice can get migrationpolicies
```bash
spicedb relationship check \
  resourcetype:migrationpolicies get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can list migrationpolicies
```bash
spicedb relationship check \
  resourcetype:migrationpolicies list user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can watch migrationpolicies
```bash
spicedb relationship check \
  resourcetype:migrationpolicies watch user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice CANNOT delete migrationpolicies (read-only)
```bash
spicedb relationship check \
  resourcetype:migrationpolicies delete user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_DENIED` (migrationpolicies is read-only)

#### Test: Alice CANNOT create migrationpolicies (read-only)
```bash
spicedb relationship check \
  resourcetype:migrationpolicies create user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_DENIED` (migrationpolicies is read-only)

### 6. VMI Subresources - GET permissions

#### Test: Alice can get VMI console
```bash
spicedb relationship check \
  subresource:vmi-console get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI VNC
```bash
spicedb relationship check \
  subresource:vmi-vnc get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI VNC screenshot
```bash
spicedb relationship check \
  subresource:vmi-vnc-screenshot get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI portforward
```bash
spicedb relationship check \
  subresource:vmi-portforward get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI guestosinfo
```bash
spicedb relationship check \
  subresource:vmi-guestosinfo get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI filesystemlist
```bash
spicedb relationship check \
  subresource:vmi-filesystemlist get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI userlist
```bash
spicedb relationship check \
  subresource:vmi-userlist get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI SEV fetchcertchain
```bash
spicedb relationship check \
  subresource:vmi-sev-fetchcertchain get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI SEV querylaunchmeasurement
```bash
spicedb relationship check \
  subresource:vmi-sev-querylaunchmeasurement get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI usbredir
```bash
spicedb relationship check \
  subresource:vmi-usbredir get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 7. VMI Subresources - UPDATE permissions

#### Test: Alice can update VMI pause
```bash
spicedb relationship check \
  subresource:vmi-pause update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI unpause
```bash
spicedb relationship check \
  subresource:vmi-unpause update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI addvolume
```bash
spicedb relationship check \
  subresource:vmi-addvolume update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI removevolume
```bash
spicedb relationship check \
  subresource:vmi-removevolume update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI freeze
```bash
spicedb relationship check \
  subresource:vmi-freeze update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI unfreeze
```bash
spicedb relationship check \
  subresource:vmi-unfreeze update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI softreboot
```bash
spicedb relationship check \
  subresource:vmi-softreboot update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI SEV setupsession
```bash
spicedb relationship check \
  subresource:vmi-sev-setupsession update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VMI SEV injectlaunchsecret
```bash
spicedb relationship check \
  subresource:vmi-sev-injectlaunchsecret update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 8. VM Subresources - GET permissions

#### Test: Alice can get VM expand-spec
```bash
spicedb relationship check \
  subresource:vm-expand-spec get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VM portforward
```bash
spicedb relationship check \
  subresource:vm-portforward get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 9. VM Subresources - UPDATE permissions

#### Test: Alice can update VM start
```bash
spicedb relationship check \
  subresource:vm-start update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM stop
```bash
spicedb relationship check \
  subresource:vm-stop update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM restart
```bash
spicedb relationship check \
  subresource:vm-restart update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM addvolume
```bash
spicedb relationship check \
  subresource:vm-addvolume update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM removevolume
```bash
spicedb relationship check \
  subresource:vm-removevolume update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM migrate
```bash
spicedb relationship check \
  subresource:vm-migrate update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update VM memorydump
```bash
spicedb relationship check \
  subresource:vm-memorydump update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 10. Cluster-scoped subresource

#### Test: Alice can update expand-vm-spec
```bash
spicedb relationship check \
  resourcetype:expand-vm-spec update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 11. Resource Instance Permissions

#### Test: Alice can get specific VM instance
```bash
spicedb relationship check \
  resource:vm-test-1 get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can delete specific VM instance
```bash
spicedb relationship check \
  resource:vm-test-1 delete user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can update specific VMI instance
```bash
spicedb relationship check \
  resource:vmi-test-1 update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 12. Subresource Instance Permissions

#### Test: Alice can update VM start on specific instance
```bash
spicedb relationship check \
  subresource_instance:vm-test-1-start update user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

#### Test: Alice can get VMI console on specific instance
```bash
spicedb relationship check \
  subresource_instance:vmi-test-1-console get user:alice \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_ALLOWED`

### 13. Negative Tests - User without role

#### Test: Bob CANNOT get virtualmachines (no role assignment)
```bash
spicedb relationship check \
  resourcetype:virtualmachines get user:bob \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_DENIED`

#### Test: Bob CANNOT delete virtualmachineinstances
```bash
spicedb relationship check \
  resourcetype:virtualmachineinstances delete user:bob \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_DENIED`

#### Test: Bob CANNOT update VM start subresource
```bash
spicedb relationship check \
  subresource:vm-start update user:bob \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```
**Expected**: `PERMISSION_DENIED`

## Listing Permissions

### List all permissions alice has on virtualmachines
```bash
spicedb relationship lookup \
  user:alice resourcetype:virtualmachines \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```

### List all users who can delete virtualmachines
```bash
spicedb relationship lookup-subjects \
  resourcetype:virtualmachines delete user \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```

### List all resources alice can access
```bash
spicedb relationship lookup-resources \
  user:alice resourcetype \
  --endpoint localhost:50051 \
  --token "somerandomkeyhere"
```

## Batch Testing Script

You can create a script to run all tests:

```bash
#!/bin/bash

ENDPOINT="localhost:50051"
TOKEN="somerandomkeyhere"

# Test function
test_permission() {
  local resource=$1
  local permission=$2
  local user=$3
  local expected=$4

  result=$(spicedb relationship check \
    "$resource" "$permission" "$user" \
    --endpoint "$ENDPOINT" \
    --token "$TOKEN" 2>&1)

  if [[ "$result" == *"$expected"* ]]; then
    echo "✓ PASS: $user $permission $resource"
  else
    echo "✗ FAIL: $user $permission $resource (expected: $expected, got: $result)"
  fi
}

# Run tests
echo "Testing ClusterRole RBAC Schema..."
echo ""

echo "=== Basic Resource Type Permissions ==="
test_permission "resourcetype:virtualmachines" "get" "user:alice" "PERMISSION_ALLOWED"
test_permission "resourcetype:virtualmachines" "delete" "user:alice" "PERMISSION_ALLOWED"
test_permission "resourcetype:virtualmachineinstances" "create" "user:alice" "PERMISSION_ALLOWED"

echo ""
echo "=== Read-only Permissions ==="
test_permission "resourcetype:migrationpolicies" "get" "user:alice" "PERMISSION_ALLOWED"
test_permission "resourcetype:migrationpolicies" "delete" "user:alice" "PERMISSION_DENIED"

echo ""
echo "=== Subresource Permissions ==="
test_permission "subresource:vm-start" "update" "user:alice" "PERMISSION_ALLOWED"
test_permission "subresource:vmi-console" "get" "user:alice" "PERMISSION_ALLOWED"

echo ""
echo "=== Negative Tests ==="
test_permission "resourcetype:virtualmachines" "get" "user:bob" "PERMISSION_DENIED"
test_permission "resourcetype:virtualmachineinstances" "delete" "user:bob" "PERMISSION_DENIED"

echo ""
echo "Tests completed!"
```

Save this as `run-tests.sh`, make it executable (`chmod +x run-tests.sh`), and run it to verify all permissions.
