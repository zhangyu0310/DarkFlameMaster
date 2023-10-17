package main

import (
	"DarkFlameMaster/tools/manager/config"
	"DarkFlameMaster/tools/manager/model"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

var (
	configPath = *flag.String("config", "./configure.json", "config path")
)

func main() {
	if err := config.ReadConfig(configPath); err != nil {
		os.Exit(1)
	}
	m := model.NewModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
