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

	tocIndex := findTableOfContents(content)
	if tocIndex == -1 {
		fmt.Printf("Warning: Table of contents non trovato in %s, processo tutto il file\n", path)
		tocIndex = 0
	}

	header := content[:tocIndex]
	body := content[tocIndex:]

	tagged := make(map[int]bool)

	type match struct {
		start int
		end   int
		term  string
	}
	var matches []match

	for _, term := range terms {
		escapedTerm := regexp.QuoteMeta(term)

		pattern := fmt.Sprintf(`\b%s\b`, escapedTerm)
		re := regexp.MustCompile(pattern)

		indices := re.FindAllStringIndex(body, -1)
		for _, loc := range indices {

			if !isInsideGloss(body, loc[0]) {
				matches = append(matches, match{
					start: loc[0],
					end:   loc[1],
					term:  term,
				})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})

	var validMatches []match
	for i := 0; i < len(matches); i++ {
		current := matches[i]

		if tagged[current.start] {
			continue
		}

		hasOverlap := false
		for j := i + 1; j < len(matches) && matches[j].start < current.end; j++ {

			if len(matches[j].term) > len(current.term) {
				hasOverlap = true
				break
			}
		}

		if !hasOverlap {
			validMatches = append(validMatches, current)

			for pos := current.start; pos < current.end; pos++ {
				tagged[pos] = true
			}
		}
	}

	for i := len(validMatches) - 1; i >= 0; i-- {
		m := validMatches[i]
		body = body[:m.start] + fmt.Sprintf("#gloss[%s]", m.term) + body[m.end:]
	}

	content = header + body

	if content != originalContent {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Errore scrittura %s: %v\n", path, err)
		} else {
			fmt.Printf("Modificato: %s\n", path)
		}
	}
}

func findTableOfContents(content string) int {
	patterns := []string{
		`#outline\(\)`,
		`#outline\[`,
		`= Indice`,
		`= Table of Contents`,
		`== Indice`,
	}

	minIndex := -1

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		loc := re.FindStringIndex(content)
		if loc != nil {
			// Trova la fine della riga dopo il pattern
			endOfLine := strings.Index(content[loc[1]:], "\n")
			if endOfLine != -1 {
				index := loc[1] + endOfLine + 1
				if minIndex == -1 || index < minIndex {
					minIndex = index
				}
			}
		}
	}

	return minIndex
}

func isInsideGloss(content string, pos int) bool {
	start := pos - 1
	depth := 0

	for start >= 0 {
		if content[start] == ']' {
			depth++
		} else if content[start] == '[' {
			depth--
			if depth < 0 {
				if start >= 6 && content[start-6:start] == "#gloss" {
					return true
				}
				return false
			}
		}
		start--
	}

	return false
}
