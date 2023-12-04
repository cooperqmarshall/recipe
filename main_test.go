package main

import (
	"testing"
)

func TestParse(t *testing.T) {
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
		"3/4 cup milk"}

	for i, ingredient := range r.Ingredients {
		if ingredients[i] != ingredient {
			t.Errorf("Recipe Ingredient incorrect. Expected: %s, Actual: %s", ingredients[i], ingredient)
		}
	}

	instructions := []string{
    "Preheat the oven to 350 degrees F. Grease and flour a 9x9 inch pan.",
		"In a large bowl, combine flour, sugar, baking powder, and salt.",
		"Mix in the butter, eggs, and milk.",
		"Spread into the prepared pan.",
		"Bake for 30 to 35 minutes, or until firm.",
		"Allow to cool and enjoy."}

	for i, instruction := range r.Instructions {
		if instructions[i] != instruction  {
			t.Errorf("Recipe Instruction incorrect. Expected: %s, Actual: %s", instructions[i], instruction)
		}
	}
}
