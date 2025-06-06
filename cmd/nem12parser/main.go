package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/lixuaniscas/nem12parser/internal/parser"
)

func main() {
    inputFile := flag.String("input", "", "Path to the input NEM12 file")
    outputFile := flag.String("output", "inserts.sql", "Path to the output SQL file")
    flag.Parse()

    if *inputFile == "" {
        log.Fatal("Input file is required")
    }

    inserts, malformed := parser.ParseFile(*inputFile)

    err := os.WriteFile(*outputFile, []byte(inserts), 0644)
    if err != nil {
        log.Fatalf("Error writing output: %v", err)
    }

    fmt.Printf("Parsing complete. SQL written to %s. Malformed lines: %d\n", *outputFile, malformed)
}
