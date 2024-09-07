package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ajtroup1/GoDoc/internal/generator"
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
	task := flag.String("task", "", Red+"Specify the task to run (e.g., save, gen)"+Reset)
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
		fmt.Println(Red + "Unknown task. Please specify 'save' or 'gen'." + Reset)
	}
	fmt.Print("" + Reset) // Prevents compiler error if all fmt are disabled

}

func save(settings *utils.SettingManager) {
	// Parse the src code into the heirarchal structure
	parser := parser.New(settings.Settings)
	parser.ParseProject()

	// If there were parsing errors, log them to the user
	if len(parser.Errors) > 0 {
		for _, err := range parser.Errors {
			log.Printf(Red+"Parsing error: %v\n"+Reset, err)
		}
	}
}

func gen(settings *utils.SettingManager) {
	// Generate the program's documentation from the heirarchal structure
	generator := generator.New(settings.Settings)
	generator.GenerateDocs()

	// If there were generation errors, log them to the user
	// In the end application, there should never be genereation errors
	if len(generator.Errors) > 0 {
		for _, err := range generator.Errors {
			log.Printf(Red+"Generation error: %v\n"+Reset, err)
		}
	}
}
