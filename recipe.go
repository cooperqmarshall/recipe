package recipe

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"
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
	Image              []string
}

type Recipe struct {
	Name         string
	Ingredients  []string
	Instructions []string
	ImageUrl     string
	ThumbnailUrl string
	jsonld       recipeJsonld
}

// Read recipe information from jsonld blob into recipe properties
func (r *Recipe) Read_jsonld(b []byte) error {
	err := json.Unmarshal(b, &r.jsonld)
	if err != nil {
		if err.Error() != "json: cannot unmarshal object into Go struct field recipeJsonld.@graph.Image of type []string" {
			return err
		}
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
			err1 := r.parse_instructions(i.RecipeInstructions)
			err2 := r.parse_images(i.Image)
			return errors.Join(err1, err2)
		}
	}

	return errors.New("Unable to find full recipe in ldjson")
}

// Parses a []string containing urls to images. Stores the url to the full
// image in recipe.ImageUrl and a smaller image in recipe.ThumbnailUrl
func (r *Recipe) parse_images(images []string) error {
	if len(images) == 0 {
		return errors.New("no images found")
	}

	shortest_url := images[0]
	// matches two numbers with an "x" in between e.g. "400x300"
	re := regexp.MustCompile(`.*?(\d+)x(\d+).*`)
	for _, image_url := range images {
		if len(image_url) < len(shortest_url) {
			shortest_url = image_url
		}

		matches := re.FindStringSubmatch(image_url)
		if len(matches) > 1 {
			width, err := strconv.ParseFloat(matches[1], 0)
			if err != nil {
				continue
			}
			height, err := strconv.ParseFloat(matches[2], 0)
			if err != nil {
				continue
			}
            // want the 4x3 aspect ratio image for the thumbnail
			if width/height == 4.0/3.0 {
				r.ThumbnailUrl = image_url
			}
		}
	}
    // assuming shortest url is the one without the aspect ratio defined and
    // will be the highest resolution option
	r.ImageUrl = shortest_url
	return nil
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
