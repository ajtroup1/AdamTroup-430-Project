package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ajtroup1/GoDoc/internal/parser"
	"github.com/ajtroup1/GoDoc/utils"
)

const (
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)

func main() {
	task := flag.String("task", "", "Specify the task to run (e.g., save, gen)")
	flag.Parse()

	// Retreive settings
	settings, err := utils.NewSettings()
	if err != nil {
		log.Fatalf(Red+"Error reading settings: %v"+Reset, err)
	}

	// Execute the appropriate task based on the flag
	switch *task {
	case "save":
		save(settings)
	case "gen":
		gen(settings)
	default:
		fmt.Println("Unknown task. Please specify 'save' or 'gen'.")
	}
	fmt.Print("" + Reset) // Prevents compiler error if all fmt are disabled

}

func save(settings *utils.SettingManager) {
	parser := parser.New(settings.Settings)
	parser.ParseProject()
}

func gen(settings *utils.SettingManager) {

}
