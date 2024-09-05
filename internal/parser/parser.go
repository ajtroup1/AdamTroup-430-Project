package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajtroup1/GoDoc/internal/models"
)

type Parser struct {
	settings     models.Settings
	src          string
	position     int
	readPosition int
	ch           byte
	Packages     []models.Package
	Errors       []error
}

func New(settings models.Settings) *Parser {
	return &Parser{settings: settings}
}

func (p *Parser) ParseProject() {
	// walk all files in the project path and read their src code
	rootDir := p.settings.ProjectPath

	// Walk through all the files in the directory
	_ = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			p.Errors = append(p.Errors, fmt.Errorf("error accessing path %s: %v", path, err))
			return err
		}

		// Process only regular files (ignore directories)
		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			err := p.readFile(path)
			if err != nil {
				p.Errors = append(p.Errors, fmt.Errorf("error reading file %s: %v", path, err))
			}
		}
		return nil
	})
}

func (p *Parser) readFile(filePath string) error {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Process the file content (in this case, just storing it in src)
	p.src = string(content)

	// You can add further logic here to parse the source code
	fmt.Printf("Reading file: %s\n", filePath)
	// Example: You can invoke a method here to parse the content
	p.parseSourceCode()

	return nil
}

func (p *Parser) parseSourceCode() {
    fmt.Printf("\n%s\n", p.src)
}

func (p *Parser) readChar() {
	if p.readPosition >= len(p.src) { // Check if the end of input is reached
		p.ch = 0 // Null character indicating end of input
	} else {
		p.ch = p.src[p.readPosition] // Read the current character
	}
	p.position = p.readPosition // Update the current position
	p.readPosition += 1         // Move to the next character
}

func (p *Parser) peekChar(ahead int) byte {
	if p.readPosition+ahead >= len(p.src) {
		return 0
	}
	return p.src[p.readPosition+ahead]
}

func (p *Parser) skipWhitespace() {
	for p.ch == ' ' || p.ch == '\t' || p.ch == '\n' || p.ch == '\r' {
		p.readChar()
	}
}
