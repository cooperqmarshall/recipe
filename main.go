package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type Step struct {
	Type            string `json:"@type"`
	Text            string
	ItemListElement []Step
}

type RecipeJsonld struct {
	Name               string
	RecipeIngredient   []string
	RecipeInstructions []Step
}

type Recipe struct {
	Name         string
	Ingredients  []string
	Instructions []string
	jsonld       RecipeJsonld
}

// Read recipe information from jsonld blob into recipe properties
func (r *Recipe) read_jsonld(b []byte) error {
	err := json.Unmarshal(b, &r.jsonld)
	if err != nil {
		return err
	}

	r.Name = r.jsonld.Name
	r.Ingredients = r.jsonld.RecipeIngredient
	r.parse_instructions(r.jsonld.RecipeInstructions)

	return err
}

// Extracts the instruction steps from recipeInstructions into []string.
// This handels both HowToStep or HowToSection elements
func (r *Recipe) parse_instructions(steps []Step) error {
	for _, s := range steps {
		// HowToStep type will have a non-nil Text attribute
		if s.Type == "HowToStep" {
			if s.Text == "" {
				return errors.New("HowtoStep does not contain 'Text' key")
			}
			r.Instructions = append(r.Instructions, s.Text)
		} else if s.Type == "HowToSection" {
			if s.ItemListElement == nil {
				return errors.New("HowToSection does not contain 'ItemListElement' key")
			}
			// HowToSection type will have a non-nil ItemListElement containing []HowToStep
			r.parse_instructions(s.ItemListElement)
		} else {
			return errors.New("Unexpected Step type")
		}
	}

	return nil
}

func get_stdin() []byte {
	std_in_reader := bufio.NewReader(os.Stdin)
	std_in_text, err := std_in_reader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	return std_in_text
}

func main() {
}
