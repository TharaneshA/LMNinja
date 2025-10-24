package attacks

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/fs"
	"math/rand"
	"path"
	"strings"
	"time"

	"embed"
)

//go:embed all:hex_phi
var hexPhiFS embed.FS

type AttackCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AttackPrompt struct {
	Category string
	Prompt   string
}


var categoryNameMap = map[string]string{
	"category_1.csv":  "Illegal Activity",
	"category_3.csv":  "Hate / Harass / Violence",
	"category_4.csv":  "Malware",
	"category_5.csv":  "Physical Harm",
	"category_6.csv":  "Economic Harm",
	"category_7.csv":  "Adult Content",
	"category_8.csv":  "Fraud Deception",
	"category_9.csv":  "Political Campaigning",
	"category_10.csv": "Privacy Violation Activity",
	"category_11.csv": "Tailored Financial Advice",
}

var (
	cachedCategories []AttackCategory
	cachedPrompts    []AttackPrompt
	isCacheBuilt     bool
)

func buildCache() error {
	if isCacheBuilt {
		return nil
	}

	categories := []AttackCategory{}
	prompts := []AttackPrompt{}

	files, err := fs.ReadDir(hexPhiFS, "hex_phi")
	if err != nil {
		return fmt.Errorf("could not read embedded hex_phi directory: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && strings.HasSuffix(fileName, ".csv") {
			// Get the category name from our reliable map
			categoryName, ok := categoryNameMap[fileName]
			if !ok {
				continue // Skip files not in our map (like the removed category_2.csv)
			}

			filePath := path.Join("hex_phi", fileName)
			content, err := hexPhiFS.ReadFile(filePath)
			if err != nil {
				continue
			}

			categories = append(categories, AttackCategory{ID: fileName, Name: categoryName})

			// Read the CSV data (skipping the header comment line)
			firstLineEnd := bytes.IndexByte(content, '\n')
			if firstLineEnd == -1 {
				firstLineEnd = 0 // Handle files with no newline
			}
			
			reader := csv.NewReader(bytes.NewReader(content[firstLineEnd+1:]))
			records, err := reader.ReadAll()
			if err != nil {
				continue
			}

			for _, record := range records {
				if len(record) > 0 && record[0] != "" {
					prompts = append(prompts, AttackPrompt{
						Category: fileName,
						Prompt:   record[0],
					})
				}
			}
		}
	}
	
	cachedCategories = categories
	cachedPrompts = prompts
	isCacheBuilt = true
	
	return nil
}

func GetCategories() ([]AttackCategory, error) {
	if err := buildCache(); err != nil {
		return nil, err
	}
	return cachedCategories, nil
}

func GetPrompts(categoryIDs []string, limit int) ([]string, error) {
	if err := buildCache(); err != nil {
		return nil, err
	}
	var filteredPrompts []string
	categorySet := make(map[string]bool)
	for _, id := range categoryIDs {
		categorySet[id] = true
	}
	for _, p := range cachedPrompts {
		if _, ok := categorySet[p.Category]; ok {
			filteredPrompts = append(filteredPrompts, p.Prompt)
		}
	}
	if limit <= 0 || limit >= len(filteredPrompts) {
		return filteredPrompts, nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(filteredPrompts), func(i, j int) {
		filteredPrompts[i], filteredPrompts[j] = filteredPrompts[j], filteredPrompts[i]
	})
	return filteredPrompts[:limit], nil
}