package parser

import (
    "strings"
    "testing"
    "time"
)

func TestProcessFileBasic(t *testing.T) {
    input := `
100,NEM12,200506081149,UNITEDDP,NEMMCO
200,NEM1201009,E1E2,1,E1,N1,01009,kWh,30,20050610
300,20050301,0,0,0,0,0,0,0,0,0,0,0,0,0.461,0.810
300,20050302,0,0,0,0,0,0,0,0,0,0,0,0,0.235,0.567
900
`
    r := strings.NewReader(input)
    cfg := Config{BatchSize: 2, Workers: 2}
    res, err := ProcessFile(r, cfg)
    if err != nil {
        t.Fatalf("ProcessFile failed: %v", err)
    }

    if res.TotalValid != 2 {
        t.Errorf("Expected 2 valid readings, got %d", res.TotalValid)
    }
    if res.MalformedCount != 0 {
        t.Errorf("Expected 0 malformed lines, got %d", res.MalformedCount)
    }
    if len(res.SQLBatches) == 0 {
        t.Errorf("Expected some SQL batches, got 0")
    }
}

func TestBuildSQLInserts(t *testing.T) {
    nmi := "NEM1201009"
    date := time.Date(2005, 3, 1, 0, 0, 0, 0, time.UTC)
    consumptions := []float64{0.1, 0.2, 0.3, 0.4}
    batches := buildSQLInserts(nmi, date, consumptions, 30, 2)

    if len(batches) != 2 {
        t.Errorf("Expected 2 batches, got %d", len(batches))
    }
    expectedPrefix := "INSERT INTO meter_readings"
    for _, batch := range batches {
        if !strings.HasPrefix(batch, expectedPrefix) {
            t.Errorf("Batch missing prefix: %s", batch)
        }
    }
}
