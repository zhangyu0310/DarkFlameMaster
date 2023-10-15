package main

import (
	"DarkFlameMaster/tools/manager/model"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	m := model.NewModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
