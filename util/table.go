package util

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintTable(t [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for _, row := range t {
		// Append an empty column so that the last real column is right-aligned too.
		r := make([]string, len(row), len(row)+1)
		copy(r, row)
		r = append(r, "")
		fmt.Fprintln(w, strings.Join(r, "\t"))
	}
	w.Flush()
}
