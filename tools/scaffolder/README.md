# Atlas Controller Scaffolder

A tool to generate Kubernetes controllers for MongoDB Atlas resources based on CRD configurations.

> [!WARNING]
> This is the experimental tool. Bugs and issues are still present

## Overview

This scaffolder generates Kubernetes controllers that follow the MongoDB Atlas Kubernetes operator patterns, including:

- **State machine-based controllers** with proper lifecycle management
- **Translation layers** for Atlas SDK integration
- **Service interfaces** with appropriate Atlas API mappings
- **License headers** and proper package structure

## Dependencies

### Required Local Repositories

The following repositories must be available locally in the parent directory:

```
work/
├── atlas-controller-scaffolder/    # This tool
├── atlas2crd/                      # Config definitions
└── mongodb-atlas-kubernetes/       # Target operator
```

### Setup

1. **Generate CRD types**
    Use `crd2go` tool from the [github.com/crd2go/crd2go](https://github.com/crd2go/crd2go) to generate go types for CRDs. In the AKO repository is installed a a go tool:

    ```bash
    go tool crd2go -input=./pkg/crd2go/samples/crds.yaml -output=../atlas-controller-scaffolder/pkg/api/v1
    ```

2. **Build the scaffolder:**
   ```bash
   cd ../scaffolder
   go build -o ./bin/scaffolder ./cmd/main.go
   ```


## Usage

### List Available CRDs

```bash
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml --list
```

### Generate Controller

```bash
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml --crd <CRD_KIND>
```

**Examples:**

```bash
# Generate Team controller
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml --crd Team

# Generate Organization controller
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml --crd Organization

# Generate DatabaseUser controller
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml --crd DatabaseUser
```

### Show Available CRDs 

If you don't specify a CRD, the tool will show you all available options:

```bash
./bin/ako-controller-scaffolder --config ../atlas2crd/config.yaml
```

## Generated Structure

The tool generates controllers in the MongoDB Atlas Kubernetes operator directory:

```
../mongodb-atlas-kubernetes/internal/
├── controller/{crd_name}/
│   ├── {crd_name}_controller.go    # Main controller with reconciler
│   └── handler.go                  # State handlers and setup
└── translation/{crd_name}/
    ├── {crd_name}.go              # Translation types
    └── service.go                 # Atlas SDK service interface
```

Generated controllers support all specified CRD versions. Translation layers stabs are generated per CRD version, using the Atlas SDK specified in the atlas2crd config file

## Features

- **Dynamic API Mapping** - Automatically selects correct Atlas SDK API
- **State Machine Pattern** - Follows existing controller patterns
- **Translation Layer** - Converts between Kubernetes and Atlas types
- **License Headers** - Proper MongoDB license in all files
- **RBAC Annotations** - Kubebuilder RBAC markers included
- **Package Structure** - Consistent with existing codebase

## Available CRDs

Run with `--list` to see all available CRDs, including:

- Team, Organization, DatabaseUser
- Cluster, FlexCluster, SearchIndex
- BackupCompliancePolicy, DataFederation
- NetworkPeeringConnection, CustomRole
- And many more...

## Development

The scaffolder uses:

- **[Jennifer](https://github.com/dave/jennifer)** for Go code generation
- **YAML parsing** for atlas2crd config files
- **Atlas SDK mapping** for correct API selection
