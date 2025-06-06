package parser

import (
    "bufio"
    "fmt"
    "io"
    "strconv"
    "strings"
    "sync"
    "time"
)

type Config struct {
    BatchSize int
    Workers   int
}

type Reading struct {
    NMI         string
    Timestamp   time.Time
    Consumption float64
}

type Result struct {
    TotalValid      int
    MalformedCount  int
    MalformedSamples []string
    SQLBatches      []string
}

type workerInput struct {
    nmi         string
    intervalLen int
    intervalDate time.Time
    consumptions []float64
}

func ProcessFile(r io.Reader, cfg Config) (*Result, error) {
    scanner := bufio.NewScanner(r)
    scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024) // up to 10MB buffer for long lines

    var (
        currentNMI       string
        currentIntervalLen int
        malformedLines    []string
        malformedCount    int
        totalValid        int
        inputsCh          = make(chan workerInput, cfg.Workers*2)
        resultsCh         = make(chan []string, cfg.Workers*2)
        wgWorkers         sync.WaitGroup
    )

    // Start worker pool
    for i := 0; i < cfg.Workers; i++ {
        wgWorkers.Add(1)
        go func() {
            defer wgWorkers.Done()
            for input := range inputsCh {
                batchInserts := buildSQLInserts(input.nmi, input.intervalDate, input.consumptions, input.intervalLen, cfg.BatchSize)
                resultsCh <- batchInserts
            }
        }()
    }

    // Producer reads lines
    var sqlBatches []string
    go func() {
        for batch := range resultsCh {
            sqlBatches = append(sqlBatches, batch...)
        }
    }()

    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Split(line, ",")

        if len(fields) < 1 {
            malformedLines = append(malformedLines, line)
            malformedCount++
            continue
        }

        switch fields[0] {
        case "200":
            // Example 200 line:
            // 200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
            if len(fields) < 10 {
                malformedLines = append(malformedLines, line)
                malformedCount++
                continue
            }
            currentNMI = fields[1]
            intervalLen, err := strconv.Atoi(fields[8])
            if err != nil {
                malformedLines = append(malformedLines, line)
                malformedCount++
                continue
            }
            currentIntervalLen = intervalLen

        case "300":
            // 300 line format:
            // 300,20050301,0,0,0,0,...consumptions...
            if currentNMI == "" || currentIntervalLen == 0 {
                // no active 200 record context yet
                malformedLines = append(malformedLines, line)
                malformedCount++
                continue
            }
            if len(fields) < 15 { // minimal consumptions expected, but can vary
                malformedLines = append(malformedLines, line)
                malformedCount++
                continue
            }
            dateStr := fields[1]
            intervalDate, err := time.Parse("20060102", dateStr)
            if err != nil {
                malformedLines = append(malformedLines, line)
                malformedCount++
                continue
            }

            consumptions := []float64{}
            // consumption values start from field index 14 (0-based)
            for _, cStr := range fields[14:] {
                cStr = strings.TrimSpace(cStr)
                if cStr == "" || cStr == "A" {
                    consumptions = append(consumptions, 0.0)
                    continue
                }
                v, err := strconv.ParseFloat(cStr, 64)
                if err != nil {
                    consumptions = append(consumptions, 0.0) // or skip, or count malformed?
                    continue
                }
                consumptions = append(consumptions, v)
            }

            inputsCh <- workerInput{
                nmi: currentNMI,
                intervalLen: currentIntervalLen,
                intervalDate: intervalDate,
                consumptions: consumptions,
            }

            totalValid++

        default:
            // ignore other record types like 100, 500, 900 etc
        }
    }

    close(inputsCh)
    wgWorkers.Wait()
    close(resultsCh)

    return &Result{
        TotalValid: totalValid,
        MalformedCount: malformedCount,
        MalformedSamples: malformedLines,
        SQLBatches: sqlBatches,
    }, nil
}

func buildSQLInserts(nmi string, date time.Time, consumptions []float64, intervalLen int, batchSize int) []string {
    inserts := []string{}
    var values []string

    for i, consumption := range consumptions {
        // Calculate timestamp for each interval:
        ts := date.Add(time.Duration(i*intervalLen) * time.Minute).Format("2006-01-02 15:04:05")
        val := fmt.Sprintf("('%s','%s',%.6f)", nmi, ts, consumption)
        values = append(values, val)

        if len(values) >= batchSize {
            stmt := fmt.Sprintf("INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES %s;", strings.Join(values, ","))
            inserts = append(inserts, stmt)
            values = []string{}
        }
    }
    if len(values) > 0 {
        stmt := fmt.Sprintf("INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES %s;", strings.Join(values, ","))
        inserts = append(inserts, stmt)
    }
    return inserts
}
