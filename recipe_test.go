package recipe

import (
	"testing"
)

func TestFullParse(t *testing.T) {
	// tests name, ingredients, and []HowToStep instructions
	b := []byte(`{
    "@context": "https://schema.org/",
    "@type": "Recipe",
    "name": "Party Coffee Cake",
    "image": [
      "https://example.com/photos/1x1/photo.jpg",
      "https://example.com/photos/4x3/photo.jpg",
      "https://example.com/photos/16x9/photo.jpg"
    ],
    "author": {
      "@type": "Person",
      "name": "Mary Stone"
    },
    "datePublished": "2018-03-10",
    "description": "This coffee cake is awesome and perfect for parties.",
    "prepTime": "PT20M",
    "cookTime": "PT30M",
    "totalTime": "PT50M",
    "keywords": "cake for a party, coffee",
    "recipeYield": "10",
    "recipeCategory": "Dessert",
    "recipeCuisine": "American",
    "nutrition": {
      "@type": "NutritionInformation",
      "calories": "270 calories"
    },
    "recipeIngredient": [
      "2 cups of flour",
      "3/4 cup white sugar",
      "2 teaspoons baking powder",
      "1/2 teaspoon salt",
      "1/2 cup butter",
      "2 eggs",
      "3/4 cup milk"
      ],
    "recipeInstructions": [
      {
        "@type": "HowToStep",
        "name": "Preheat",
        "text": "Preheat the oven to 350 degrees F. Grease and flour a 9x9 inch pan.",
        "url": "https://example.com/party-coffee-cake#step1",
        "image": "https://example.com/photos/party-coffee-cake/step1.jpg"
      },
      {
        "@type": "HowToStep",
        "name": "Mix dry ingredients",
        "text": "In a large bowl, combine flour, sugar, baking powder, and salt.",
        "url": "https://example.com/party-coffee-cake#step2",
        "image": "https://example.com/photos/party-coffee-cake/step2.jpg"
      },
      {
        "@type": "HowToStep",
        "name": "Add wet ingredients",
        "text": "Mix in the butter, eggs, and milk.",
        "url": "https://example.com/party-coffee-cake#step3",
        "image": "https://example.com/photos/party-coffee-cake/step3.jpg"
      },
      {
        "@type": "HowToStep",
        "name": "Spread into pan",
        "text": "Spread into the prepared pan.",
        "url": "https://example.com/party-coffee-cake#step4",
        "image": "https://example.com/photos/party-coffee-cake/step4.jpg"
      },
      {
        "@type": "HowToStep",
        "name": "Bake",
        "text": "Bake for 30 to 35 minutes, or until firm.",
        "url": "https://example.com/party-coffee-cake#step5",
        "image": "https://example.com/photos/party-coffee-cake/step5.jpg"
      },
      {
        "@type": "HowToStep",
        "name": "Enjoy",
        "text": "Allow to cool and enjoy.",
        "url": "https://example.com/party-coffee-cake#step6",
        "image": "https://example.com/photos/party-coffee-cake/step6.jpg"
      }
    ]
  }`)

	r := Recipe{}
	err := r.read_jsonld(b)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r.Name != "Party Coffee Cake" {
		t.Errorf("Recipe Name incorrect. Expected: Party Coffee Cake, Actual: %s", r.Name)
	}

	ingredients := []string{
		"2 cups of flour",
		"3/4 cup white sugar",
		"2 teaspoons baking powder",
		"1/2 teaspoon salt",
		"1/2 cup butter",
		"2 eggs",
		"3/4 cup milk",
	}
	ingredientsMatch(t, ingredients, r)

	instructions := []string{
		"Preheat the oven to 350 degrees F. Grease and flour a 9x9 inch pan.",
		"In a large bowl, combine flour, sugar, baking powder, and salt.",
		"Mix in the butter, eggs, and milk.",
		"Spread into the prepared pan.",
		"Bake for 30 to 35 minutes, or until firm.",
		"Allow to cool and enjoy.",
	}
	instructionsMatch(t, instructions, r)
}

// tests instruction parsing when HowToSection type
func TestSimpleInstructions(t *testing.T) {
	b := []byte(`{"recipeInstructions": 
        [{
          "@type": "HowToSection",
          "name": "Assemble the pie",
          "itemListElement": [
            {
              "@type": "HowToStep",
              "text": "In large bowl, gently mix filling ingredients; spoon into crust-lined pie plate."
            }, {
              "@type": "HowToStep",
              "text": "Top with second crust. Cut slits or shapes in several places in top crust."
            }
          ]
        }]
    }`)

	r := Recipe{}
	err := r.read_jsonld(b)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r.Name != "" {
		t.Errorf("Recipe Name incorrect. Expected: \"\", Actual: %s", r.Name)
	}

	if r.Ingredients != nil {
		t.Errorf("Recipe Ingredient incorrect. Expected: nil, Actual: %s", r.Ingredients)
	}

	instructions := []string{
		"In large bowl, gently mix filling ingredients; spoon into crust-lined pie plate.",
		"Top with second crust. Cut slits or shapes in several places in top crust.",
	}
	instructionsMatch(t, instructions, r)
}

// test if HowToStep is missing
func TestBadHowToSection(t *testing.T) {
	b := []byte(`{"recipeInstructions": 
        [{
          "@type": "HowToSection",
          "name": "Assemble the pie",
        }]
    }`)

	r := Recipe{}
	err := r.read_jsonld(b)
	if err == nil {
		t.Errorf("Expected to Error due to missing 'itemListElement' in HowToSection")
	}
}

// test if HowToStep is missing
func TestBadHowToStep(t *testing.T) {
	b := []byte(`{"recipeInstructions": 
        [{
          "@type": "HowToStep",
        }]
    }`)

	r := Recipe{}
	err := r.read_jsonld(b)
	if err == nil {
		t.Errorf("Expected to Error due to missing 'Text' in HowToStep")
	}
}

func instructionsMatch(t *testing.T, instructions []string, r Recipe) {
	if len(instructions) != len(r.Instructions) {
		t.Errorf("Recipe Instruction incorrect length. Expected: %d, Actual: %d", len(instructions), len(r.Instructions))
	} else {
		for i, instruction := range instructions {
			if r.Instructions[i] != instruction {
				t.Errorf("Recipe Instruction incorrect. Expected: %s, Actual: %s", instructions[i], instruction)
			}
		}
	}
}

func ingredientsMatch(t *testing.T, ingredients []string, r Recipe) {
	if len(ingredients) != len(r.Ingredients) {
		t.Errorf("Recipe Ingrediets incorrect length. Expected: %d, Actual: %d", len(ingredients), len(r.Ingredients))
	} else {
		for i, ingredient := range ingredients {
			if r.Ingredients[i] != ingredient {
				t.Errorf("Recipe Ingredient incorrect. Expected: %s, Actual: %s", ingredients[i], ingredient)
			}
		}
	}
}
