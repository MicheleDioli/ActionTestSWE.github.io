package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	GlossaryFile = "src/glossary.typ"
)

var TargetFolders = []string{
	"src/RTB",
	"src/Candidatura",
}

var FilesToIgnore = map[string]bool{
	"LetteraPresentazione.typ": true,
}

func main() {
	if _, err := os.Stat(GlossaryFile); os.IsNotExist(err) {
		fmt.Printf("No glossario in: %s\n", GlossaryFile)
		return
	}

	fmt.Println("🔍 Estrazione termini da", GlossaryFile, "...")
	terms, err := extractTerms(GlossaryFile)
	if err != nil {
		fmt.Printf("Lett glossario: %v\n", err)
		return
	}

	for _, folder := range TargetFolders {

		err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("Errore in %s: %v\n", path, err)
				return filepath.SkipDir
			}

			if d.IsDir() {
				return nil
			}

			fileName := filepath.Base(path)

			if filepath.Ext(path) == ".typ" && !FilesToIgnore[fileName] {

				if filepath.Clean(path) == filepath.Clean(GlossaryFile) {
					return nil
				}

				processFile(path, terms)
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Errore scansione %s: %v\n", folder, err)
		}
	}

	fmt.Println("\nFINE")
}

func extractTerms(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`term\s*:\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	var terms []string
	for _, match := range matches {
		if len(match) > 1 {
			terms = append(terms, match[1])
		}
	}

	sort.Slice(terms, func(i, j int) bool {
		return len(terms[i]) > len(terms[j])
	})

	return terms, nil
}

func processFile(path string, terms []string) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Impossibile leggere %s: %v\n", path, err)
		return
	}
	content := string(contentBytes)
	originalContent := content

	for _, term := range terms {
		escapedTerm := regexp.QuoteMeta(term)
		pattern := fmt.Sprintf(`(#gloss\[)?\b(%s)\b`, escapedTerm)
		re := regexp.MustCompile(pattern)

		content = re.ReplaceAllStringFunc(content, func(match string) string {
			if strings.HasPrefix(match, "#gloss[") {
				return match
			}
			return fmt.Sprintf("#gloss[%s]", term)
		})
	}

	if content != originalContent {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Errore scrittura %s: %v\n", path, err)
		} else {
			fmt.Printf("Modificato: %s\n", path)
		}
	}
}
