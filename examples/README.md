# Terraform Import Remover Examples

This directory contains example Terraform configurations that demonstrate the functionality of the Terraform Import Remover tool.

## Directory Structure

```
examples/
├── main.tf                    # Basic example with import blocks
├── modules/                   # Nested directory with modules
│   ├── main.tf                # Module configuration with import block
│   └── networking/            # Nested module
│       └── vpc.tf             # VPC configuration with import blocks
└── edge_cases/                # Examples of edge cases
    ├── empty.tf               # Empty file
    ├── only_import.tf         # File with only import blocks
    └── commented.tf           # File with commented import blocks
```

## Usage

1. Run the tool on the examples directory:

```bash
./terraform-import-remover ./examples
```

2. Observe the output statistics showing files processed and import blocks removed.

3. Check the modified files to see that the import blocks have been removed.

## What to Expect

- The tool will process all `.tf` files recursively in the examples directory.
- It will remove all `import` blocks from these files.
- It will report statistics about the files processed and blocks removed.
- The files will maintain their structure and formatting, with only the `import` blocks removed.

## Edge Cases

The `edge_cases` directory demonstrates how the tool handles various special cases:

- `empty.tf`: An empty file (tool should process it without errors)
- `only_import.tf`: A file containing only `import` blocks (will become empty after processing)
- `commented.tf`: A file with commented out `import` blocks (comments should remain untouched)
