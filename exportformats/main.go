package main

import "fmt"

type Exporter interface {
	Export() string
}

type CSVExport struct {
	Filename string
}

func (c CSVExport) Export() string {
	return fmt.Sprintf("CSV: %s.csv", c.Filename)
}

type JSONExport struct {
	Filename string
}

func (j JSONExport) Export() string {
	return fmt.Sprintf("JSON: %s.json", j.Filename)
}

func RunExport(e Exporter) {
	fmt.Println(e.Export())
}

func main() {
	csv := CSVExport{Filename: "movies"}
	json := JSONExport{Filename: "movies"}

	RunExport(csv)
	RunExport(json)
}
