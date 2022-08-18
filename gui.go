package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	clipboard "github.com/atotto/clipboard"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	goutils "github.com/simonski/goutils"
)

type GUI struct {
	DB *KPDB
}

func NewGUI(db *KPDB) *GUI {
	g := GUI{DB: db}
	return &g
}

func (g *GUI) Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	list := widgets.NewList()
	list.Title = "Entries"
	entries := g.DB.GetEntriesSortedByUpdatedThenKey()
	s := make([]string, 0)
	max_key := 0
	for _, e := range entries {
		s = append(s, e.Key)
		max_key = goutils.Max(max_key, len(e.Key))
	}
	list.Rows = s
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false
	list.SetRect(0, 0, 25, 25)

	table_width := 80

	table := widgets.NewTable()
	table.Title = "Entries"
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.TextAlignment = ui.AlignLeft
	table.RowSeparator = false
	table.FillRow = false
	table.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(26, 0, 26+table_width, 25)

	table.ColumnWidths = []int{max_key + 1, 80}

	index := 0
	current_entry := entries[index]
	g.drawTable(table, current_entry, index, table_width)

	ui.Render(list)
	ui.Render(table)

	output := ""
	quit := false
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				index, _ = strconv.Atoi(e.ID)
				index -= 1
			case "g", "G":
				index = len(entries) - 1
			case "s", "j", "<Down>":
				index += 1
			case "w", "k", "<Up>":
				index -= 1
			case "<Return>", "<Enter>":
				entry, _ := g.DB.GetDecrypted(current_entry.Key)
				clipboard.WriteAll(entry.Value)
				quit = true
			case "q", "Q", "<C-c>":
				clipboard.WriteAll("quit")
				output = "quit"
				quit = true
			}

			if index >= len(entries) {
				index = 0
			} else if index < 0 {
				index = len(entries) - 1
			}
			current_entry = entries[index]
			list.ScrollTop()
			list.ScrollAmount(index)
		}

		if quit {
			break
		}
		g.drawTable(table, current_entry, index, table_width)
		ui.Render(list)
		ui.Render(table)
	}
	ui.Close()
	if output != "" {
		fmt.Println(output)
	}
}

func (g *GUI) drawTable(table *widgets.Table, entry DBEntry, index int, table_width int) {

	table.Title = fmt.Sprintf("Table [%v]", index)
	table.Rows = [][]string{}

	rows := [][]string{}
	rows = append(rows, []string{"Key", entry.Key})
	splits := g.SplitText(entry.Description, table_width)
	if len(splits) > 1 {
		for index := 0; index < len(splits); index++ {
			if index == 0 {
				rows = append(rows, []string{"Description", splits[index]})
			} else {
				rows = append(rows, []string{"", splits[index]})
			}
		}
	} else {
		rows = append(rows, []string{"Description", entry.Description})
	}
	rows = append(rows, []string{"", ""})

	rows = append(rows, []string{"Type", entry.Type})
	rows = append(rows, []string{"Url", entry.Url})
	rows = append(rows, []string{"Username", entry.Username})

	rows = append(rows, []string{"LastUpdated", entry.LastUpdated.Format(time.RFC822)})
	rows = append(rows, []string{"Created", entry.LastUpdated.Format(time.RFC822)})

	splits = g.SplitText(entry.Notes, table_width)
	if len(splits) > 1 {
		for index := 0; index < len(splits); index++ {
			if index == 0 {
				rows = append(rows, []string{"Notes", splits[index]})
			} else {
				rows = append(rows, []string{"", splits[index]})
			}
		}
	} else {
		rows = append(rows, []string{"Notes", entry.Notes})
	}
	rows = append(rows, []string{"", ""})

	// 	[]string{"Created", entry.LastUpdated.Format(time.RFC822)},
	// 	[]string{"Notes", entry.Notes},
	// 	[]string{"Url", entry.Url},
	// 	[]string{"Type", entry.Type},
	// 	[]string{"Value", entry.Value},
	// }

	table.Rows = rows

}

// splits a text if it exceeds a maximum width or has newlines.
func (g *GUI) SplitText(s string, width int) []string {
	splits := strings.Split(s, "\\n")
	result := make([]string, 0)
	for index := 0; index < len(splits); index++ {
		if len(splits[index]) <= width {
			result = append(result, splits[index])
		} else {
			remaining := splits[index]
			for {
				if len(remaining) > width {
					left := remaining[0:width]
					remaining = remaining[width:]
					result = append(result, left)
				} else {
					result = append(result, remaining)
					break
				}
			}
		}
	}
	return result

}
