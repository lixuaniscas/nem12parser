# nem12parser

## Getting Started

### Prerequisites
- Go 1.20 or higher installed.  
  Download from https://golang.org/dl/

### Installation

Clone this repository:

```bash
git clone https://github.com/lixuaniscas/nem12parser.git
cd nem12parser
```

Build the executable:

```bash
go build -o nem12parser ./cmd/nem12parser
```

### Usage

Run the parser on your `.nem12` input file and output SQL inserts:

```bash
./nem12parser -input=testdata/sample.nem12 -batch=1000 -workers=4

```

#### Command-line flags:

- `-input` (required): Path to the input NEM12 file.
- `-output` (optional): Path to output SQL file (defaults to stdout).
- `-batch` (optional): Batch size for SQL inserts (default 5000).

---

## Testing

Run unit tests with verbose output:

```bash
go test ./internal/parser -v
```

---

## Project Structure

```
nem12parser/
├── cmd/
│   └── nem12parser/        # Main executable
│       └── main.go
├── internal/
│   └── parser/             # Core parsing and SQL generation logic
│       ├── parser.go
│       └── parser_test.go
├── sample/
│   └── sample.nem12        # Sample input file
├── go.mod
└── README.md
```

---

## Design Decisions & Rationale

- **Standard Library Only:** Avoids dependency issues and eases deployment.
- **Streaming & Concurrency:** Uses streaming file reading and worker pools to handle large files efficiently without high memory usage.
- **Batch Inserts:** Generates batched SQL insert statements to optimize database ingestion.
- **Validation & Reporting:** Captures malformed lines to maintain data quality.
- **Unit Tests & Logging:** Improves maintainability and debugging.

---

## Future Improvements

- Support for multiple output formats (e.g., CSV, JSON).
- Database direct insert support with transaction handling.
- Configurable concurrency and chunk size.
- More extensive schema validation and error recovery.
