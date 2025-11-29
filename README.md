# toxjc

**toxjc** - File format conversion tool (to XLS/JSON/CSV)

Convert between CSV, JSON, XLSX, and XLS file formats with automatic type detection.

> The name "toxjc" comes from: **to** **X**(LS) **J**(SON) **C**(SV)

## Installation

### Method 1: Homebrew (Recommended)

```bash
brew install yejune/tap/toxjc
```

### Method 2: Using go install

```bash
go install github.com/yejune/toxjc@latest
toxjc install
```

## Usage

### Convert files

```bash
toxjc <input> <output>
```

The output format is determined by the file extension.

### Examples

```bash
# Excel to CSV
toxjc data.xlsx output.csv

# CSV to JSON
toxjc data.csv output.json

# JSON to Excel
toxjc data.json output.xlsx

# Old Excel format
toxjc legacy.xls output.csv
```

### Detect file type

```bash
toxjc detect <file>
```

```bash
$ toxjc detect unknown.dat
csv
```

### Version

```bash
toxjc version
```

## Supported Formats

| Format | Input | Output |
|--------|-------|--------|
| CSV    | ✅    | ✅     |
| XLSX   | ✅    | ✅     |
| XLS    | ✅    | ❌     |
| JSON   | ✅    | ✅     |

## Features

- Automatic file type detection
- Encoding auto-detection (UTF-8, UTF-16, EUC-KR)
- Handles BOM (Byte Order Mark)
- Simple CLI interface

## License

MIT License
