package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type HowToStep struct {
	Text string
}

type HowToSection struct {
	ItemListElement []HowToStep
	Text            string
}

type RecipeJsonld struct {
	Name               string         `json:"name"`
	RecipeIngredient   []string       `json:"recipeIngredient"`
	RecipeInstructions []HowToSection `json:"recipeInstructions"`
}

type Recipe struct {
	Name         string
	Ingredients  []string
	Instructions []string
	jsonld       RecipeJsonld
}

// read recipe information from jsonld blob into Recipe properties
func (r *Recipe) read_jsonld(b []byte) error {
	var rld RecipeJsonld
	err := json.Unmarshal(b, &rld)

	r.Name = rld.Name
	r.Ingredients = rld.RecipeIngredient
	r.Instructions = rld.parse_instructions()

	return err
}

func (rld *RecipeJsonld) parse_instructions() []string {
	instructions := []string{}

	for _, s := range rld.RecipeInstructions {
		// if HowToStep
		if s.Text != "" {
			instructions = append(instructions, s.Text)
		} else if s.ItemListElement != nil {
			// if CreativeWork
			parse_howtosteps(s.ItemListElement, instructions)
		}
	}

	return instructions
}

func parse_howtosteps(steps []HowToStep, instructions []string) {
	for _, step := range steps {
		if step.Text != "" {
			instructions = append(instructions, step.Text)
		}
	}
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
