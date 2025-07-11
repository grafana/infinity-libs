package csvframer

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/gframer"
	"github.com/grafana/infinity-libs/lib/go/jsonframer"
)

type FramerOptions struct {
	FrameName          string
	Columns            []gframer.ColumnSelector
	Delimiter          string
	SkipLinesWithError bool
	Comment            string
	RelaxColumnCount   bool
	NoHeaders          bool
	FramerType         jsonframer.FramerType // `gjson` | `jsonata` | `jq`
	RootSelector       string
}

func ToFrame(csvString string, options FramerOptions) (frame *data.Frame, err error) {
	if strings.TrimSpace(csvString) == "" {
		return frame, ErrEmptyCsv
	}
	r := csv.NewReader(strings.NewReader(csvString))
	r.LazyQuotes = true
	if options.Comment != "" {
		r.Comment = rune(options.Comment[0])
	}
	if options.Delimiter != "" {
		r.Comma = rune(options.Delimiter[0])
	}
	if options.RelaxColumnCount {
		r.FieldsPerRecord = -1
	}
	parsedCSV := [][]string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err == nil {
			parsedCSV = append(parsedCSV, record)
			continue
		}
		if !options.SkipLinesWithError {
			return frame, errors.Join(ErrReadingCsvResponse, fmt.Errorf("%w, %v", err, record))
		}
	}
	out := []interface{}{}
	header := []string{}
	records := [][]string{}
	if !options.NoHeaders {
		header = parsedCSV[0]
		for idx, hItem := range header {
			for _, col := range options.Columns {
				if col.Selector == hItem && col.Alias != "" && options.RootSelector == "" {
					header[idx] = col.Alias
				}
			}
		}
		records = parsedCSV[1:]
	}
	if options.NoHeaders {
		records = parsedCSV
		if len(records) > 0 {
			for i := 0; i < len(records[0]); i++ {
				header = append(header, fmt.Sprintf("%d", i+1))
			}
		}
	}
	for _, row := range records {
		item := map[string]interface{}{}
		for colId, col := range header {
			if colId < len(row) {
				item[col] = row[colId]
			}
		}
		out = append(out, item)
	}
	framerOptions := gframer.FramerOptions{FrameName: options.FrameName, Columns: options.Columns}
	if options.RootSelector != "" {
		outObj, err := ApplyRootSelector(out, options.RootSelector, options.FramerType)
		if err != nil {
			return nil, err
		}
		return gframer.ToDataFrame(outObj, framerOptions)
	}
	return gframer.ToDataFrame(out, framerOptions)
}

func ApplyRootSelector(csvArray []any, rootSelector string, framerType jsonframer.FramerType) (any, error) {
	outStringBytes, err := json.Marshal(csvArray)
	if err != nil {
		return nil, err
	}
	outString, err := jsonframer.ApplyRootSelector(string(outStringBytes), rootSelector, framerType)
	if err != nil {
		return nil, err
	}
	var outObj any
	jsonUnMarshallErr := json.Unmarshal([]byte(outString), &outObj)
	if jsonUnMarshallErr != nil {
		return nil, err
	}
	return outObj, nil
}
