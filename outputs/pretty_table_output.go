package output

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// PrettyTableOutput is ...
type PrettyTableOutput struct {
	options         *PrettyTableOutputOptions
	writer          *tablewriter.Table
	firstRow        []string
	header          []string
	minOutputLength int
}

// PrettyTableOutputOptions is ...
type PrettyTableOutputOptions struct {
	WriteHeader  bool
	WriteTo      io.Writer
	RepeatHeader bool
}

// NewPrettyTableOutput is ...
func NewPrettyTableOutput(opts *PrettyTableOutputOptions) *PrettyTableOutput {
	prettyTableOutput := &PrettyTableOutput{
		options: opts,
		writer:  tablewriter.NewWriter(opts.WriteTo),
	}

	prettyTableOutput.writer.SetAutoWrapText(false)
	if prettyTableOutput.options.RepeatHeader {
		prettyTableOutput.writer.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	}

	return prettyTableOutput
}

// Show is ...
func (o *PrettyTableOutput) Show(rows *sql.Rows) {
	cols, colsErr := rows.Columns()

	if colsErr != nil {
		log.Fatalln(colsErr)
	}

	o.writer.SetColWidth(128)
	if o.options.WriteHeader {
		o.writer.SetHeader(cols)
		o.writer.SetHeaderLine(true)
		o.writer.SetAutoFormatHeaders(false)
	}

	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols))

	for i := range cols {
		dest[i] = &rawResult[i]
	}

	j := 0
	for rows.Next() {
		rows.Scan(dest...)

		for i, raw := range rawResult {

			if _, err := strconv.ParseInt(string(raw), 10, 64); err == nil {
				result[i] = string(raw) + ".00"
			} else if f, err := strconv.ParseFloat(string(raw), 64); err == nil {
				result[i] = fmt.Sprintf("%.2f", f)
			} else {
				result[i] = string(raw)
			}

		}

		o.writer.Append(result)
		j++
		if o.options.RepeatHeader {
			if j%30 == 0 {
				o.writer.Render()
				o.writer.ClearRows()
			}
		}
	}

	if o.options.RepeatHeader {
		o.writer.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	}
	o.writer.Render()
	rows.Close()
}
