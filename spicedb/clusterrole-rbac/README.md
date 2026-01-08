# ClusterRole RBAC Schema for SpiceDB

This directory contains a SpiceDB schema design for modeling Kubernetes ClusterRole permissions, specifically for the `kubevirt.io:admin` ClusterRole.

## Overview

The schema models the following Kubernetes RBAC concepts:

- **ClusterRole**: Cluster-wide roles that can be assigned to users
- **API Groups**: Logical groupings of Kubernetes resources (e.g., `kubevirt.io`, `snapshot.kubevirt.io`)
- **Resource Types**: Types of Kubernetes resources (e.g., `virtualmachines`, `virtualmachineinstances`)
- **Resources**: Specific instances of resource types
- **Subresources**: Special endpoints on resources (e.g., `/start`, `/stop`, `/console`, `/vnc`)
- **Verbs/Permissions**: Actions that can be performed (get, list, watch, create, update, patch, delete, deletecollection)

## Schema Design

### Entities

1. **user**: Represents a subject that can be assigned roles
2. **clusterrole**: A cluster-wide role with permissions
3. **apigroup**: A Kubernetes API group
4. **resourcetype**: A type of Kubernetes resource within an API group
5. **resource**: A specific instance of a resource type
6. **subresource**: A subresource path (e.g., `/console`, `/start`)
7. **subresource_instance**: A specific instance of a subresource

### Permission Model

The schema uses a hierarchical permission model:

```
ClusterRole
    ├── Assigned to Users
    ├── Admin of API Groups
    └── Has permissions on Resource Types
            ├── get, list, watch, create, update, patch, delete, deletecollection
            └── Subresources with get/update permissions
```

### Key Features

1. **Role-based access**: Users are assigned to ClusterRoles, which grant permissions
2. **Hierarchical permissions**: API groups can have parent groups, and permissions inherit
3. **Granular verb control**: Each verb (get, list, watch, etc.) is modeled separately
4. **Subresource support**: Models Kubernetes subresources like `/start`, `/stop`, `/console`
5. **Instance-level permissions**: Can model permissions on specific resource instances

## Files

- `schema.yaml`: SpiceDB schema definition
- `relationships.yaml`: Sample relationships for the `kubevirt.io:admin` ClusterRole
- `test-queries.md`: Example queries to test the schema

## Example Queries

### Check if a user can get a virtualmachine

```bash
spicedb relationship check \
  resourcetype:virtualmachines get user:alice
```

Expected: `PERMISSION_ALLOWED` (alice has the kubevirt-admin role)

### Check if a user can delete a virtualmachineinstance

```bash
spicedb relationship check \
  resourcetype:virtualmachineinstances delete user:alice
```

Expected: `PERMISSION_ALLOWED`

### Check if a user can update a VM subresource (start)

```bash
spicedb relationship check \
  subresource:vm-start update user:alice
```

Expected: `PERMISSION_ALLOWED`

### Check if a user can get VMI console

```bash
spicedb relationship check \
  subresource:vmi-console get user:alice
```

Expected: `PERMISSION_ALLOWED`

### Check permissions for a user without the role

```bash
spicedb relationship check \
  resourcetype:virtualmachines get user:bob
```

Expected: `PERMISSION_DENIED` (bob is not assigned to kubevirt-admin)

### Check read-only permission on migrationpolicies

```bash
spicedb relationship check \
  resourcetype:migrationpolicies get user:alice
```

Expected: `PERMISSION_ALLOWED`

```bash
spicedb relationship check \
  resourcetype:migrationpolicies delete user:alice
```

Expected: `PERMISSION_DENIED` (migrationpolicies only has get, list, watch - no delete)

## Mapping from ClusterRole YAML

The original ClusterRole has these rule groups:

### 1. VMI Subresources - GET only
- `virtualmachineinstances/console`
- `virtualmachineinstances/vnc`
- `virtualmachineinstances/vnc/screenshot`
- `virtualmachineinstances/portforward`
- `virtualmachineinstances/guestosinfo`
- `virtualmachineinstances/filesystemlist`
- `virtualmachineinstances/userlist`
- `virtualmachineinstances/sev/fetchcertchain`
- `virtualmachineinstances/sev/querylaunchmeasurement`
- `virtualmachineinstances/usbredir`

Mapped to subresources with `get_role` relation.

### 2. VMI Subresources - UPDATE only
- `virtualmachineinstances/pause`
- `virtualmachineinstances/unpause`
- `virtualmachineinstances/addvolume`
- `virtualmachineinstances/removevolume`
- `virtualmachineinstances/freeze`
- `virtualmachineinstances/unfreeze`
- `virtualmachineinstances/softreboot`
- `virtualmachineinstances/sev/setupsession`
- `virtualmachineinstances/sev/injectlaunchsecret`

Mapped to subresources with `update_role` relation.

### 3. VM Subresources - GET only
- `virtualmachines/expand-spec`
- `virtualmachines/portforward`

Mapped to subresources with `get_role` relation.

### 4. VM Subresources - UPDATE only
- `virtualmachines/start`
- `virtualmachines/stop`
- `virtualmachines/restart`
- `virtualmachines/addvolume`
- `virtualmachines/removevolume`
- `virtualmachines/migrate`
- `virtualmachines/memorydump`

Mapped to subresources with `update_role` relation.

### 5. Full CRUD on kubevirt.io resources
All verbs (get, delete, create, update, patch, list, watch, deletecollection) on:
- `virtualmachines`
- `virtualmachineinstances`
- `virtualmachineinstancepresets`
- `virtualmachineinstancereplicasets`
- `virtualmachineinstancemigrations`

Each verb mapped to a separate role relation (e.g., `get_role`, `delete_role`).

### 6. Full CRUD on other API groups
Same pattern for:
- `snapshot.kubevirt.io`: virtualmachinesnapshots, virtualmachinesnapshotcontents, virtualmachinerestores
- `export.kubevirt.io`: virtualmachineexports
- `clone.kubevirt.io`: virtualmachineclones
- `instancetype.kubevirt.io`: virtualmachineinstancetypes, virtualmachineclusterinstancetypes, etc.
- `pool.kubevirt.io`: virtualmachinepools

### 7. Read-only on migrationpolicies
Only get, list, watch (no create, update, delete) on:
- `migrations.kubevirt.io/migrationpolicies`

## Testing the Schema

1. Start SpiceDB with the schema:
```bash
spicedb serve --grpc-preshared-key "somerandomkeyhere" &
```

2. Load the schema:
```bash
spicedb schema write --schema schema.yaml
```

3. Load the relationships:
```bash
cat relationships.yaml | spicedb relationship import
```

4. Run test queries (see examples above)

## Design Considerations

### Why separate resourcetype and resource?

This allows modeling permissions at both the type level (all VMs) and instance level (specific VM). Most K8s RBAC operates at the type level, but SpiceDB can support instance-level permissions if needed.

### Why model verbs as separate relations?

This provides maximum flexibility. Different roles might have different combinations of verbs. For example:
- A "reader" role might only have `get_role`, `list_role`, `watch_role`
- An "admin" role has all verb roles
- A "deployer" role might have `create_role`, `update_role`, but not `delete_role`

### Why use apigroup->access in resourcetype permissions?

This creates a fallback mechanism where having admin access to an API group grants access to all its resources. This matches the hierarchical nature of Kubernetes RBAC.

### How to handle cluster-scoped vs namespaced resources?

The current schema models cluster-scoped resources (ClusterRole). For namespaced resources, you would add a `namespace` definition and relation to resources, and filter permissions by namespace membership.

## Future Enhancements

1. **Namespace support**: Add namespace entities for namespaced resources
2. **RoleBindings**: Model the binding between subjects (users, groups, service accounts) and roles
3. **Groups**: Add group entity and group membership
4. **Service Accounts**: Model service accounts as subjects
5. **Conditional permissions**: Use caveats for time-based or attribute-based access control
6. **Audit trail**: Track who has access to what and why (via relationship traversal)
