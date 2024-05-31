package recipe

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type step struct {
	Type            string `json:"@type"`
	Text            string
	ItemListElement []step
}

type recipeJsonld struct {
	Name               string
	RecipeIngredient   []string
	RecipeInstructions []step
	Graph              []recipeJsonld `json:"@graph"`
}

type ldjson struct {
}

type Recipe struct {
	Name         string
	Ingredients  []string
	Instructions []string
	jsonld       recipeJsonld
}

// Read recipe information from jsonld blob into recipe properties
func (r *Recipe) Read_jsonld(b []byte) error {
	err := json.Unmarshal(b, &r.jsonld)
	if err != nil {
		return err
	}

	if r.jsonld.RecipeIngredient != nil && r.jsonld.RecipeInstructions != nil {
		r.Name = r.jsonld.Name
		r.Ingredients = r.jsonld.RecipeIngredient
		err = r.parse_instructions(r.jsonld.RecipeInstructions)
		return err
	}

	// search through graph array for recipe ldjsons
	for _, i := range r.jsonld.Graph {
		if i.RecipeIngredient != nil && i.RecipeInstructions != nil {
			r.Name = i.Name
			r.Ingredients = i.RecipeIngredient
			err = r.parse_instructions(i.RecipeInstructions)
			return err
		}
	}

	return errors.New("Unable to find full recipe in ldjson")
}

// Extracts the instruction steps from recipeInstructions into []string.
// This handels both HowToStep or HowToSection elements
func (r *Recipe) parse_instructions(steps []step) error {
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
			err := r.parse_instructions(s.ItemListElement)
			if err != nil {
				return err
			}
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
