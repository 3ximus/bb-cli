package util

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Table(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("\033[1;30m─\033[m")
	// table.SetNoWhiteSpace(true)
	// table.SetTablePadding("\t")
	headerColors := make([]tablewriter.Colors, len(header))
	for i := range headerColors {
		headerColors[i] = tablewriter.Colors{tablewriter.FgHiBlueColor, tablewriter.Bold}
	}
	table.SetHeaderColor(headerColors...)
	table.AppendBulk(data)

	os.Stdout.Write([]byte("\n"))
	table.Render()
	os.Stdout.Write([]byte("\n"))

}
