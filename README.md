Getting Started
Prerequisites
Go 1.20 or higher installed. Download from https://golang.org/dl/

Installation
Clone this repository:

bash
Copy
Edit
git clone https://github.com/yourusername/nem12parser.git
cd nem12parser
Build the executable:

bash
Copy
Edit
go build -o nem12parser ./cmd/nem12parser
Usage
Run the parser on your .nem12 input file and output SQL inserts:

bash
Copy
Edit
./nem12parser -input=path/to/sample.nem12 -output=output.sql
Command-line flags:

-input (required): path to the input NEM12 file.

-output (optional): path to output SQL file (defaults to stdout).

-batch (optional): batch size for SQL inserts (default 5000).

Testing
Run unit tests:

bash
Copy
Edit
go test ./internal/parser -v
Project Structure
pgsql
Copy
Edit
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
Design Decisions & Rationale
Standard Library Only: Avoids dependency issues and eases deployment.

Streaming & Concurrency: Uses streaming file read and worker pools to scale for large files without memory bloat.

Batch Inserts: Generates batched SQL insert statements to optimize DB ingestion.

Validation & Reporting: Captures malformed lines to maintain data quality.

Unit Tests & Logging: Improves maintainability and debugging.

Future Improvements
Support for multiple output formats (e.g., CSV, JSON).

Database direct insert support with transaction handling.

Configurable concurrency and chunk size.

More extensive schema validation and error recovery.
