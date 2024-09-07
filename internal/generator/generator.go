package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ajtroup1/GoDoc/internal/models"
)

type Generator struct {
	Settings models.Settings
	Packages []models.Package
	Errors   []error
}

func New(settings models.Settings) *Generator {
	return &Generator{Settings: settings}
}

func (g *Generator) GenerateDocs() {
	g.readJSON()

	if g.Settings.DocGenFormat == "markdown" {
		// Create or open the markdown file
		docPath := "./Docs.md"
		if g.Settings.ProjectName != "" {
			docPath = fmt.Sprintf("%s/%s.md", g.Settings.DocGenPath, g.Settings.ProjectName)
		} else {
			if g.Settings.DocGenPath != "./" {
				docPath = fmt.Sprintf("%s/Docs.md", g.Settings.DocGenPath)
			}
		}
		fmt.Printf("%s\n", docPath)
		file, err := os.Create(docPath)
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("failed to open/create documentation path from settings '%s', ensure it exists", g.Settings.DocGenPath))
			return
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		g.generateHeaderMD(writer)
		// g.generateTOCMD()
		g.generateBodyMD(writer)
	}
}

func (g *Generator) readJSON() {
	// Open the JSON file
	file, err := os.Open("./godoc_output.json")
	if err != nil {
		g.Errors = append(g.Errors, err)
		return
	}
	defer file.Close()

	// Decode the JSON data into the Packages field
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&g.Packages)
	if err != nil {
		g.Errors = append(g.Errors, err)
		return
	}
}

func (g *Generator) generateHeaderMD(writer *bufio.Writer) {
	if g.Settings.ProjectName != "" {
		_, err := writer.WriteString(fmt.Sprintf("# %s\n", g.Settings.ProjectName))
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("error writing header to markdown: %v", err))
		}
		if g.Settings.ProjectDesc != "" {
			_, err = writer.WriteString(fmt.Sprintf("%s\n", g.Settings.ProjectDesc))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing header to markdown: %v", err))
			}
		} else {
			_, err = writer.WriteString("\n")
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing header to markdown: %v", err))
			}
		}
	} else {
		_, err := writer.WriteString("# GoDoc generator documentation\n")
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("error writing header to markdown: %v", err))
		}
	}
}

func (g *Generator) generateBodyMD(writer *bufio.Writer) {
	if len(g.Packages) == 0 {
		g.Errors = append(g.Errors, fmt.Errorf("no packages found in the stored comment tree"))
		return
	}
	_, err := writer.WriteString("## Packages:\n")
	if err != nil {
		g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
		return
	}
	for _, pkg := range g.Packages {
		_, err := writer.WriteString(fmt.Sprintf("  - ### Package: `%s`\n", pkg.Name))
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
			return
		}
		if pkg.Desc != "" {
			_, err = writer.WriteString(fmt.Sprintf("    %s\n\n", pkg.Desc))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
				return
			}
		}
		if len(pkg.Files) > 0 {
			g.writeFile(pkg.Files, writer)
		} else {
			g.Errors = append(g.Errors, fmt.Errorf("no files in package '%s'", pkg.Name))
		}

		if len(pkg.Types) > 0 {
			_, err = writer.WriteString("      - #### Types:\n")
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}

			for _, _type := range pkg.Types {
				_, err = writer.WriteString(fmt.Sprintf("        - **%s**\n", _type.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
					return
				}
				_, err = writer.WriteString(fmt.Sprintf("          - %s\n", _type.Desc))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
					return
				}
				_, err = writer.WriteString("          - Fields:\n")
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
					return
				}
				for _, field := range _type.Fields {
					_, err = writer.WriteString(fmt.Sprintf("            - `%s`\n              - Data type: `%s`\n              - %s\n", field.Name, field.Type, field.Desc))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
				}
			}
		}

		if len(pkg.Funcs) > 0 {
			_, err = writer.WriteString(fmt.Sprintf("      - #### Functions for `%s`:\n", pkg.Name))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}

			for _, function := range pkg.Funcs {
				_, err = writer.WriteString(fmt.Sprintf("        - **%s**\n", function.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
					return
				}
				if function.Desc != "" {
					_, err = writer.WriteString(fmt.Sprintf("          - %s\n", function.Desc))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
				}
				if function.Receiver != "" {
					_, err = writer.WriteString(fmt.Sprintf("          - %s\n", function.Receiver))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
				}
				if len(function.Params) > 0 {
					_, err = writer.WriteString("          - Parameters:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
					for _, param := range function.Params {
						_, err = writer.WriteString(fmt.Sprintf("              - `%s`\n                - Data type: `%s`\n                - %s\n", param.Name, param.Type, param.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
							return
						}
					}
				}
				if len(function.Returns) > 0 {
					_, err = writer.WriteString("          - Return values:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
					for _, ret := range function.Returns {
						_, err = writer.WriteString(fmt.Sprintf("              - `%s`\n                - %s\n", ret.Paren, ret.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
							return
						}
					}
				}
				if len(function.Responses) > 0 {
					_, err = writer.WriteString("          - HTTP responses:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
					for _, res := range function.Responses {
						_, err = writer.WriteString(fmt.Sprintf("              - `%s`\n                - %s\n", res.Paren, res.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
							return
						}
					}
				}
			}
		}

		if len(pkg.Vars) > 0 {
			_, err = writer.WriteString(fmt.Sprintf("      - #### Variables for `%s`:\n", pkg.Name))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}

			for _, variable := range pkg.Vars {
				_, err = writer.WriteString(fmt.Sprintf("        - **%s**\n", variable.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
					return
				}
				if variable.Type != "" {
					_, err = writer.WriteString(fmt.Sprintf("          - Data type: `%s`\n", variable.Type))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
				}
				if variable.Desc != "" {
					_, err = writer.WriteString(fmt.Sprintf("          - %s\n", variable.Desc))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
						return
					}
				}
			}
		}

		_, err = writer.WriteString("---\n")
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("error writing package description to markdown: %v", err))
			return
		}
	}

	// Ensure all buffered content is flushed to the file
	err = writer.Flush()
	if err != nil {
		g.Errors = append(g.Errors, fmt.Errorf("error flushing writer: %v", err))
	}
}

func (g *Generator) writeFile(files []models.File, writer *bufio.Writer) {
	_, err := writer.WriteString("      - #### Files:\n")
	if err != nil {
		g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
		return
	}
	for _, file := range files {
		_, err = writer.WriteString(fmt.Sprintf("        - `%s`\n", file.Name))
		if err != nil {
			g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
			return
		}
		if file.Desc != "" {
			_, err = writer.WriteString(fmt.Sprintf("          - %s\n", file.Desc))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
		}
		if file.Author != "" {
			_, err = writer.WriteString(fmt.Sprintf("          - Authored by: **%s**\n", file.Author))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
		}
		if file.Version != "" {
			_, err = writer.WriteString(fmt.Sprintf("          - Version: **%s**\n", file.Version))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
		}
		if file.Date != "" {
			_, err = writer.WriteString(fmt.Sprintf("          - Updated on: **%s**\n", file.Date))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
		}

		if len(file.Types) > 0 {
			_, err := writer.WriteString(fmt.Sprintf("          - **Types for file `%s`**:\n", file.Name))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
			for _, _type := range file.Types {
				_, err = writer.WriteString(fmt.Sprintf("            - %s\n", _type.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
					return
				}
				if _type.Desc != "" {
					_, err = writer.WriteString(fmt.Sprintf("            - Updated on: **%s**\n", file.Date))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
				}
				if len(_type.Fields) > 0 {
					_, err = writer.WriteString("            - Fields:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
					for _, field := range _type.Fields {
						_, err = writer.WriteString(fmt.Sprintf("                - `%s:`\n", field.Name))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - Data type: `%s`\n", field.Type))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - %s\n", field.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
					}
				}
			}
		}

		if len(file.Funcs) > 0 {
			_, err := writer.WriteString(fmt.Sprintf("          - **Functions for file `%s`**:\n", file.Name))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
			for _, function := range file.Funcs {
				_, err = writer.WriteString(fmt.Sprintf("            - %s\n", function.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
					return
				}
				if function.Desc != "" {
					_, err = writer.WriteString(fmt.Sprintf("            - Updated on: **%s**\n", file.Date))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
				}
				if len(function.Params) > 0 {
					_, err = writer.WriteString("            - Parameters:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
					for _, param := range function.Params {
						_, err = writer.WriteString(fmt.Sprintf("                - `%s:`\n", param.Name))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - Data type: `%s`\n", param.Type))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - %s\n", param.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
					}
				}
				if len(function.Returns) > 0 {
					_, err = writer.WriteString("            - Returns:\n")
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
					for _, param := range function.Params {
						_, err = writer.WriteString(fmt.Sprintf("                - `%s:`\n", param.Name))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - Data type: `%s`\n", param.Type))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
						_, err = writer.WriteString(fmt.Sprintf("                  - %s\n", param.Desc))
						if err != nil {
							g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
							return
						}
					}
				}
			}
		}
		if len(file.Vars) > 0 {
			_, err := writer.WriteString(fmt.Sprintf("          - **Variables for file `%s`**:\n", file.Name))
			if err != nil {
				g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
				return
			}
			for _, variable := range file.Vars {
				_, err = writer.WriteString(fmt.Sprintf("            - %s\n", variable.Name))
				if err != nil {
					g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
					return
				}
				if variable.Type != "" {
					_, err = writer.WriteString(fmt.Sprintf("            - Data type: `%s`\n", variable.Type))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
				}
				if variable.Desc != "" {
					_, err = writer.WriteString(fmt.Sprintf("            - %s\n", variable.Desc))
					if err != nil {
						g.Errors = append(g.Errors, fmt.Errorf("error writing package name to markdown: %v", err))
						return
					}
				}
			}
		}
	}
}
