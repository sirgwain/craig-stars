package cmd

import (
	"fmt"
	"os"

	"github.com/kataras/tablewriter"
	"github.com/lensesio/tableprinter"
)

func PrintTable(caption string, data interface{}) {

	if caption != "" {
		fmt.Printf("\n%s\n\n", caption)
	}

	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = false, false, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.HeaderBgColor = tablewriter.BgBlackColor
	printer.HeaderFgColor = tablewriter.FgGreenColor

	printer.Print(data)
}
