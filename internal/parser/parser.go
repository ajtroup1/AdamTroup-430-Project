package parser

import (
	"bufio"
	"fmt"
	"log"
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
	var comments []models.Comment
	// Walk through all the files in the directory
	err := filepath.WalkDir(p.settings.ProjectPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			if p.settings.IncludeTests || (!p.settings.IncludeTests && !strings.HasSuffix(entry.Name(), "_test.go")) {
				fileComments := p.parseSourceCode(path)
				if len(fileComments) > 0 {
					log.Printf("%d comments found in %s\n", len(fileComments), path)
					comments = append(comments, fileComments...)
				}
			}
		}
		return nil
	})

	if err != nil {
		p.Errors = append(p.Errors, fmt.Errorf("error walking through project directory: %v", err))
	}

	if len(comments) != 0 {
		p.parseComments(comments)
	}
}

func (p *Parser) parseSourceCode(filePath string) []models.Comment {
	var comments []models.Comment

	log.Printf("Reading file '%s'\n", filePath)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return comments
	}
	defer file.Close()

	// Read the content of the file into e.src
	content, err := p.readFile(file)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return comments
	}
	p.src = content

	pkgName := p.extractPkgName()

	// Now proceed to extract comments
	p.readPosition = 0
	p.readChar() // Reset to the beginning of the file to start comment extraction

	for !p.isAtEnd() {
		if p.isGoDocComment() {
			comment := p.extractBlockComment(filePath, pkgName)
			if !isEmptyComment(comment) {
				comments = append(comments, comment)
			}
		} else {
			p.readChar()
		}
	}

	return comments
}

func (p *Parser) parseComments(comments []models.Comment) {
	// First, initialize the pkgs so nodes can be assigned to them
	p.initializePackages(comments)

	// Next, initialize files (to pkgs) so unexported items can be linked to them
	p.initializeFiles(comments)

	// Test print for pkgs and files
	for i, pkg := range p.Packages {
		fmt.Printf("PACKAGE %d: %s\n%d files:\n", i+1, pkg.Name, len(pkg.Files))
		for _, file := range pkg.Files {
			fmt.Printf("  - (%s) %s\n", file.Name, file.Path)
		}
	}

	for _, comment := range comments {
		// Get the header from the comment block
		headerLine := comment.Text[0]
		keyword, err := extractKeyword(headerLine)
		// fmt.Printf("%s\n", keyword)
		if err != nil {
			p.Errors = append(p.Errors, err)
		} else {
			// Only want to evalute if there is a valid
			// text := strings.Join(comment.Text[1:], "\n") // All of the comment except the header line
			// Retreive comment information according to block type
			switch keyword {
			case "FILE":
			case "TYPE", "T":
			case "FUNCTION", "FUNC":
			case "VARIABLE", "VAR", "V":
			}
		}
	}
}

func (p *Parser) initializePackages(comments []models.Comment) {
	var pkgNames []string // Prevents duplicate pkg declaration
	// Get all pkg names
	for _, comment := range comments {
		if !contains(pkgNames, comment.Package) {
			pkgNames = append(pkgNames, comment.Package)
		}
	}

	for _, comment := range comments {
		// Get the header from the comment block
		headerLine := comment.Text[0]
		keyword, err := extractKeyword(headerLine)
		// fmt.Printf("%s\n", keyword)
		if err != nil {
			p.Errors = append(p.Errors, err)
		} else {
			// Only want to evalute if there is a valid
			text := strings.Join(comment.Text[1:], "\n") // All of the comment except the header line
			// Retreive comment information for pkg types. Also check for unrecognized headers here
			// Most nodes cannot be assigned to a tree since all packages may not be accounted for yet
			switch keyword {
			case "PACKAGE", "PKG", "P":
				var pkg models.Package
				tags, err := p.extractTagData(text)
				if err != nil {
					p.Errors = append(p.Errors, err)
				} else {
					for _, tag := range tags {
						// fmt.Printf("%s: %s\n", tag.Name, tag.Content)
						if tag.Name == "package" || tag.Name == "pkg" || tag.Name == "p" || tag.Name == "name" {
							pkg.Name = tag.Content
						} else if tag.Name == "description" || tag.Name == "desc" {
							pkg.Desc = tag.Content
						} else if tag.Name == "usage" || tag.Name == "u" {
							pkg.Usage = tag.Content
						} else if tag.Name == "dependency" || tag.Name == "dep" {
							dep, err := p.extractDependency(tag.Content)
							if err != nil {
								p.Errors = append(p.Errors, err)
							} else {
								pkg.Deps = append(pkg.Deps, dep)
							}
						} else {
							p.Errors = append(p.Errors, fmt.Errorf("tag name '%s' unrecognized for pkg declaration", tag.Name))
						}
					}
				}
				// Ensure the pkg name represents an actual pkg in the src code
				found := false
				for _, name := range pkgNames {
					if pkg.Name == name {
						found = true
					}
				}
				if !found {
					p.Errors = append(p.Errors, fmt.Errorf("package name '%s' does not represent a real package in the code", pkg.Name))
				} else {
					// Ensure no duplicate pkg declarations
					found = false
					for _, p := range p.Packages {
						if pkg.Name == p.Name {
							found = true
						}
					}
					if found {
						p.Errors = append(p.Errors, fmt.Errorf("duplicate package declartion for '%s'", pkg.Name))
					} else {
						p.Packages = append(p.Packages, pkg)
					}
				}
			case "FILE":
			case "TYPE", "T":
			case "FUNCTION", "FUNC":
			case "VARIABLE", "VAR", "V":
			default:
				p.Errors = append(p.Errors, fmt.Errorf("unrecognized header '%s'", keyword))
			}
		}
	}
	// Optional test print for packages
	// for i, pkg := range p.Packages {
	// 	fmt.Printf("PACKAGE %d: %s\n%s\nDeps:\n", i+1, pkg.Name, pkg.Desc)
	// 	for _, dep := range pkg.Deps {
	// 		fmt.Printf("  - %s: %s\n", dep.Name, dep.Desc)
	// 	}
	// 	fmt.Printf("\n____________\n\n")
	// }
}

func (p *Parser) initializeFiles(comments []models.Comment) {
	for _, comment := range comments {
		// Get the header from the comment block
		headerLine := comment.Text[0]
		keyword, err := extractKeyword(headerLine)
		// fmt.Printf("%s\n", keyword)
		if err != nil {
			p.Errors = append(p.Errors, err)
		} else {
			// Only want to evalute if there is a valid
			text := strings.Join(comment.Text[1:], "\n") // All of the comment except the header line
			// Retreive comment information according to block type
			switch keyword {
			case "FILE":
				var file models.File
				file.Path = comment.File
				tags, err := p.extractTagData(text)
				if err != nil {
					p.Errors = append(p.Errors, err)
				} else {
					for _, tag := range tags {
						// fmt.Printf("%s: %s\n", tag.Name, tag.Content)
						if tag.Name == "file" || tag.Name == "f" || tag.Name == "name" {
							file.Name = tag.Content
						} else if tag.Name == "description" || tag.Name == "desc" {
							file.Desc = tag.Content
						} else if tag.Name == "author" || tag.Name == "auth" || tag.Name == "a" {
							file.Author = tag.Content
						} else if tag.Name == "version" || tag.Name == "v" {
							file.Version = tag.Content
						} else if tag.Name == "date" || tag.Name == "d" {
							file.Date = tag.Content
						} else {
							p.Errors = append(p.Errors, fmt.Errorf("tag name '%s' unrecognized for file declaration", tag.Name))
						}
					}
				}

				// Allocate the file to its pkg
				found := false
				for i := range p.Packages {
					if p.Packages[i].Name == comment.Package {
						found = true
						p.Packages[i].Files = append(p.Packages[i].Files, file)
					}
				}

				if !found {
					p.Errors = append(p.Errors, fmt.Errorf("no package found for file '%s'", file.Path))
				}
			}
		}
	}
}

func extractKeyword(line string) (string, error) {
	// Trim spaces and check if the line starts with `-- `
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "-- ") {
		keyword := ""
		// Remove the `-- ` prefix
		keyword = strings.TrimSpace(strings.TrimPrefix(line, "-- "))
		if keyword == "" {
			return "", fmt.Errorf("extracted keyword is blank")
		}
		// Further trim spaces around the keyword
		return keyword, nil
	} else {
		return "", fmt.Errorf("expected '-- ' before comment block type")
	}
}

func (p *Parser) extractTagData(text string) ([]models.Tag, error) {
	var buffer strings.Builder
	var tags []models.Tag
	name, content := "", ""
	text = strings.TrimSpace(text)
	p.src = text
	p.resetState()
	// Optionally, print the line under examination
	// fmt.Printf("%s\n", p.src)

	for p.ch != 0 {
		// fmt.Printf("%c", p.ch)
		if p.ch == '@' {
			p.readChar() // Skip @
			for p.ch != ' ' {
				buffer.WriteByte(p.ch)
				p.readChar()
			}
			p.readChar() // Skip ' '
			if len(buffer.String()) > 0 {
				name = strings.ToLower(strings.TrimSpace(buffer.String()))
			} else {
				return tags, fmt.Errorf("blank tag name")
			}
			buffer.Reset()

			// fmt.Printf("%s\n", name)

			// Tag name found, read content until another tag declaration (@)
			for p.ch != '@' && p.ch != 0 {
				buffer.WriteByte(p.ch)
				p.readChar()
			}
			if len(buffer.String()) > 0 {
				content = strings.TrimSpace(buffer.String())
				tags = append(tags, models.Tag{Name: name, Content: content})
			} else {
				return tags, fmt.Errorf("blank tag name")
			}
			buffer.Reset()
		} else {
			p.readChar()
		}
	}

	return tags, nil
}

func (p *Parser) extractDependency(content string) (models.Dependency, error) {
	var dep models.Dependency
	var buffer strings.Builder

	p.src = content
	// fmt.Printf("%s\n", p.src)

	p.resetState()

	if strings.HasPrefix(content, "(") {
		p.readChar()
		//Extract dep name
		for p.ch != 0 && p.ch != ')' {
			buffer.WriteByte(p.ch)
			p.readChar()
		}
		if len(buffer.String()) > 0 {
			dep.Name = buffer.String()
		} else {
			return dep, fmt.Errorf("blank tag name")
		}
		buffer.Reset()
		p.readChar()
		p.readChar()
		// Extract dep description
		for p.ch != 0 {
			buffer.WriteByte(p.ch)
			p.readChar()
		}
		if len(buffer.String()) > 0 {

			dep.Desc = strings.TrimSpace(buffer.String())
		} else {
			return dep, fmt.Errorf("dependency description for '%s' is blank", dep.Name)
		}
		buffer.Reset()
	} else {
		return dep, fmt.Errorf("expected '(dependencyName) Description about depencency'")
	}

	return dep, nil
}

func (p *Parser) resetState() {
	p.position = 0
	p.readPosition = 0
	p.readChar()
}

func (p *Parser) readFile(file *os.File) (string, error) {
	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (p *Parser) extractPkgName() string {
	lines := strings.Split(p.src, "\n")
	for _, line := range lines {
		// Check for the line that starts with "package "
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "package "))
		}
	}
	return ""
}

func (p *Parser) extractBlockComment(filePath, pkgName string) models.Comment {
	var lines []string
	p.advanceBy(4) // Skip /***

	for !p.isAtEnd() {
		p.skipWhitespace()

		// Check for the end of the block comment */
		if p.ch == '*' && p.peekChar(0) == '/' {
			p.advanceBy(2) // Skip closing */
			break
		}

		// Read and store each line
		var sb strings.Builder
		for p.ch != '\n' && p.ch != 0 {
			sb.WriteByte(p.ch)
			p.readChar()
		}

		line := strings.TrimSpace(sb.String())
		if len(line) > 0 {
			lines = append(lines, line)
		}

		p.readChar() // Move to the next line
	}

	if len(lines) == 0 {
		p.Errors = append(p.Errors, fmt.Errorf("no lines in file %s", filePath))
		return models.Comment{}
	}

	return models.Comment{
		File:    filePath,
		Package: pkgName,
		Text:    lines,
	}
}

func (p *Parser) advanceBy(n int) {
	for i := 0; i < n; i++ {
		p.readChar()
	}
}

func (p *Parser) isGoDocComment() bool {
	return p.ch == '/' && p.peekChar(0) == '*' && p.peekChar(1) == '*' && p.peekChar(2) == '*'
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

func isEmptyComment(comment models.Comment) bool {
	return len(comment.Text) == 0
}

func (p *Parser) isAtEnd() bool {
	return p.readPosition >= len(p.src)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
