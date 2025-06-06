package parser

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

type Record200 struct {
    NMI           string
    IntervalLength int
}

func ParseFile(filePath string) (string, int) {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var (
        current Record200
        inserts []string
        malformedCount int
    )

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "200") {
            parts := strings.Split(line, ",")
            if len(parts) >= 10 {
                current.NMI = parts[1]
                current.IntervalLength, _ = strconv.Atoi(parts[8])
            } else {
                malformedCount++
            }
        } else if strings.HasPrefix(line, "300") {
            parts := strings.Split(line, ",")
            if len(parts) < 52 {
                malformedCount++
                continue
            }

            dateStr := parts[1]
            t, err := time.Parse("20060102", dateStr)
            if err != nil {
                malformedCount++
                continue
            }

            for i := 2; i < 2+current.IntervalLength && i < len(parts); i++ {
                val := strings.TrimSpace(parts[i])
                if val == "" {
                    continue
                }
                cons, err := strconv.ParseFloat(val, 64)
                if err != nil {
                    malformedCount++
                    continue
                }
                timestamp := t.Add(time.Duration(i-2) * time.Minute * time.Duration(current.IntervalLength))
                inserts = append(inserts, fmt.Sprintf(
                    "INSERT INTO meter_readings (nmi, timestamp, consumption) VALUES ('%s', '%s', %.3f);",
                    current.NMI, timestamp.Format("2006-01-02 15:04:05"), cons,
                ))
            }
        }
    }

    return strings.Join(inserts, "\n"), malformedCount
}
