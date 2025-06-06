package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "nem12parser/internal/parser"
)

func main() {
    inputFile := flag.String("input", "", "Path to NEM12 input file")
    batchSize := flag.Int("batch", 1000, "Batch size for SQL inserts")
    workers := flag.Int("workers", 4, "Number of concurrent workers")

    flag.Parse()

    if *inputFile == "" {
        log.Fatal("Input file is required. Use -input=<path>")
    }

    f, err := os.Open(*inputFile)
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer f.Close()

    cfg := parser.Config{
        BatchSize: *batchSize,
        Workers:   *workers,
    }

    result, err := parser.ProcessFile(f, cfg)
    if err != nil {
        log.Fatalf("Processing failed: %v", err)
    }

    fmt.Printf("Parsed %d readings successfully\n", result.TotalValid)
    fmt.Printf("Malformed lines: %d\n", result.MalformedCount)
    if result.MalformedCount > 0 {
        fmt.Println("Malformed lines samples:")
        for i, line := range result.MalformedSamples {
            if i >= 10 {
                break
            }
            fmt.Println(line)
        }
    }

    // Output SQL insert statements
    fmt.Println("\n-- Generated SQL insert statements --")
    for _, batch := range result.SQLBatches {
        fmt.Println(batch)
    }
}
