package main

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Table(todo []Todo) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"S.No.", "UID", "Task", "Status", "Created At", "Completed At"})

	RoundedBorder := table.BoxStyle{
		BottomLeft:       "╰",
		BottomRight:      "╯",
		BottomSeparator:  "┴",
		EmptySeparator:   " ",
		Left:             "│",
		LeftSeparator:    "├",
		MiddleHorizontal: "─",
		MiddleSeparator:  "┼",
		MiddleVertical:   "│",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "│",
		RightSeparator:   "┤",
		TopLeft:          "╭",
		TopRight:         "╮",
		TopSeparator:     "┬",
		UnfinishedRow:    " ≈",
	}

	customStyle := table.StyleDefault
	customStyle.Box = RoundedBorder
	// t.SetStyle(customStyle)
	t.SetStyle(table.StyleRounded)

	for i, td := range todo {
		done := "❌"
		completed := "-"

		if td.Status {
			done = "✅"
			if td.CompletedAt != nil {
				completed = td.CompletedAt.Format("02-Jan-2006 3:04:05 PM")
			}
		}
		created := td.CreatedAt.Format("02-Jan-2006 3:04:05 PM")
		t.AppendRow(table.Row{
			i + 1,
			td.UID,
			td.Task,
			done,
			created,
			completed,
		})
	}
	t.Render()
}
