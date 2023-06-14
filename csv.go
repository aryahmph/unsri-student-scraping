package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVWriter struct {
	file   *os.File
	writer *csv.Writer
}

func NewCSVWriter(filename string) *CSVWriter {
	file, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)
	initHeader(writer)
	return &CSVWriter{file: file, writer: writer}
}

func (w *CSVWriter) write(data [][]string) error {
	for _, row := range data {
		err := w.writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func initHeader(writer *csv.Writer) {
	header := []string{
		"Name",
		"NIM",
		"Batch Year",
		"Image URL",
		"Gender",
		"Religion",
		"School Origin",
		"Faculty",
		"Major",
		"Major Detail",
		"Status",
	}

	err := writer.Write(header)
	if err != nil {
		panic(err)
	}
}
