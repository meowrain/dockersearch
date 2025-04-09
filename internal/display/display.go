package display

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/meowrain/dockersearch/internal/models"
	"github.com/olekukonko/tablewriter"
)

func limitString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatPullCount(count int) string {
	if count >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(count)/1000000000)
	} else if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	} else if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func DisplayPrettyTable(results *models.SearchResponse) {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n%s '%s' 的搜索结果（第 %s/%s 页，共 %s 条）\n\n",
		bold("Docker Hub"),
		bold(results.Query),
		yellow(fmt.Sprintf("%d", results.Page)),
		yellow(fmt.Sprintf("%d", results.NumPages)),
		yellow(fmt.Sprintf("%d", results.NumResults)))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "DESCRIPTION", "STARS", "OFFICIAL", "AUTOMATED", "PULLS"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, img := range results.Results {
		official := ""
		if img.IsOfficial {
			official = green("[OK]")
		}
		automated := ""
		if img.IsAutomated {
			automated = cyan("[OK]")
		}
		desc := limitString(img.Description, 60)
		stars := fmt.Sprintf("%d", img.StarCount)
		if img.StarCount > 1000 {
			stars = yellow(stars)
		}
		pulls := formatPullCount(img.PullCount)
		table.Append([]string{
			img.Name,
			desc,
			stars,
			official,
			automated,
			pulls,
		})
	}
	table.Render()
}
