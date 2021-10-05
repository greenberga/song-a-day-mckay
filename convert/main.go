package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var files = []string{
	"raw/500-451.json",
	"raw/450-401.json",
	"raw/400-351.json",
	"raw/350-301.json",
	"raw/300-251.json",
	"raw/250-201.json",
	"raw/200-151.json",
	"raw/150-101.json",
	"raw/100-51.json",
	"raw/50-1.json",
}

func readRows(f string) (rows [][]string, err error) {
	contents, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	if err := json.Unmarshal(contents, &rows); err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	return rows, nil
}

func writeCsv(rows [][]string) error {
	f, err := os.Create("songs.csv")
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	if err := csv.NewWriter(f).WriteAll(rows); err != nil {
		return fmt.Errorf("writing rows: %w", err)
	}

	return nil
}

func main() {
	var rows [][]string

	for _, f := range files {
		log.Printf("reading file: %s", f)
		next, err := readRows(f)
		if err != nil {
			log.Fatalf("failed to read rows: %v", err)
		}

		rows = append(rows, next...)
	}

	log.Printf("read %d rows", len(rows))

	if err := writeCsv(rows); err != nil {
		log.Fatalf("failed to write csv: %v", err)
	}
}
