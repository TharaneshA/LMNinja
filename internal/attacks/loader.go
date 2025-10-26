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
	"category_7.csv":  "Fraud Deception",
	"category_8.csv":  "Adult Content",
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
			categoryName, ok := categoryNameMap[fileName]
			if !ok {
				continue 
			}

			filePath := path.Join("hex_phi", fileName)
			content, err := hexPhiFS.ReadFile(filePath)
			if err != nil {
				continue
			}

			categories = append(categories, AttackCategory{ID: fileName, Name: categoryName})

			firstLineEnd := bytes.IndexByte(content, '\n')
			if firstLineEnd == -1 {
				firstLineEnd = 0 
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

func GetPrompts(categoryIDs []string, limitPerCategory int) ([]string, error) {
	if err := buildCache(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	finalPrompts := []string{}
	
	promptsByCategory := make(map[string][]string)
	for _, p := range cachedPrompts {
		promptsByCategory[p.Category] = append(promptsByCategory[p.Category], p.Prompt)
	}
	for _, categoryID := range categoryIDs {
		promptsInCat, ok := promptsByCategory[categoryID]
		if !ok {
			continue 
		}
		rand.Shuffle(len(promptsInCat), func(i, j int) {
			promptsInCat[i], promptsInCat[j] = promptsInCat[j], promptsInCat[i]
		})

		numToTake := limitPerCategory
		if limitPerCategory <= 0 || limitPerCategory > len(promptsInCat) {
			numToTake = len(promptsInCat) 
		}
		finalPrompts = append(finalPrompts, promptsInCat[:numToTake]...)
	}

	return finalPrompts, nil
}