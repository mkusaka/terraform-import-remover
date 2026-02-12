# Terraform Import Block Remover

> [!WARNING]
> This repository is kept for backward compatibility.
> Please migrate to [`tftidy`](https://github.com/mkusaka/tftidy), which unifies moved/removed/import cleanup.
>
> Migration:
> `terraform-import-remover [options] [directory]` -> `tftidy --type import [options] [directory]`

A Go tool that recursively scans Terraform files and removes all `import` blocks.

## Overview

This tool helps clean up Terraform configurations by stripping out any `import` blocks—often used to manage state imports inline but sometimes left behind after refactoring. It parses each file with Terraform’s HCL library to ensure accurate syntax handling and preserves the rest of your configuration.

## Features

- **Recursive scan**: Walks through directories to find all `.tf` files.  
- **Import block removal**: Detects and removes every `import { … }` block.  
- **In-place edits**: Modifies files directly, or simulates changes in dry-run mode.  
- **Formatting**: Runs `terraform fmt` on modified files to restore standard formatting.  
- **Verbose reporting**: Shows detailed stats on files processed, modified, and import blocks removed.  
- **Whitespace control**: Optionally normalize extra blank lines left after removal.

## Requirements

- Go 1.24 or later  
- Terraform HCL v2 (bundled via `github.com/hashicorp/hcl/v2`)

## Installation

### From Source

```bash
git clone https://github.com/mkusaka/terraform-import-remover.git
cd terraform-import-remover
go build -o terraform-import-remover cmd/terraform-import-remover/main.go
````

### Using Go Install

```bash
go install github.com/mkusaka/terraform-import-remover/cmd/terraform-import-remover@latest
```

This installs the binary to your `$GOPATH/bin`.

## Usage

```bash
terraform-import-remover [options] [directory]
```

If you omit `directory`, the current directory is used by default.

### Options

* `-help`
  Display help information and exit.
* `-version`
  Show version and exit.
* `-dry-run`
  Scan and report without writing changes.
* `-verbose`
  Print file-by-file details as you go.
* `-normalize-whitespace`
  Collapse extra blank lines after block removal (default: false).

### Example

```bash
terraform-import-remover -verbose -normalize-whitespace ./infra
```

This command will:

1. Find all `.tf` files under `./infra`.
2. Remove any `import { … }` blocks.
3. Run `terraform fmt` on updated files.
4. Print a summary of changes.

## Example Output

```
Scanning directory: ./infra
Found 20 Terraform files

Statistics:
  Files processed:        20
  Files modified:         5
  Import blocks removed:  8
  Processing time:        142.738ms
```

## How It Works

1. **Parse** each `.tf` file using the HashiCorp HCL v2 parser.
2. **Walk** the AST to locate `Block` nodes where `Type == "import"`.
3. **Remove** those blocks from the file’s AST.
4. **Serialize** the modified AST back to HCL text.
5. **Run** `terraform fmt` to clean up formatting (unless in dry-run).
6. **Report** summary statistics to the console.

## License

MIT

```
