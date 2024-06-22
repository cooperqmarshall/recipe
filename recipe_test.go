package recipe

import (
	"testing"
)

func TestFullParse(t *testing.T) {
	// tests name, ingredients, and []HowToStep instructions
	b := []byte(`
{
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
}
  `)

	r := Recipe{}
	err := r.Read_jsonld(b)
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
	b := []byte(`
{
  "recipeInstructions": [
    {
      "@type": "HowToSection",
      "name": "Assemble the pie",
      "itemListElement": [
        {
          "@type": "HowToStep",
          "text": "In large bowl, gently mix filling ingredients; spoon into crust-lined pie plate."
        },
        {
          "@type": "HowToStep",
          "text": "Top with second crust. Cut slits or shapes in several places in top crust."
        }
      ]
    }
  ],
  "recipeIngredient": []
}
    `)

	r := Recipe{}
	err := r.Read_jsonld(b)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r.Name != "" {
		t.Errorf("Recipe Name incorrect. Expected: \"\", Actual: %s", r.Name)
	}

	if len(r.Ingredients) != 0 {
		t.Errorf("Recipe Ingredient incorrect. Expected: [], Actual: %s", r.Ingredients)
	}

	instructions := []string{
		"In large bowl, gently mix filling ingredients; spoon into crust-lined pie plate.",
		"Top with second crust. Cut slits or shapes in several places in top crust.",
	}
	instructionsMatch(t, instructions, r)
}

// test if HowToStep is missing
func TestBadHowToSection(t *testing.T) {
	b := []byte(`
    {
        "recipeInstructions": 
        [{
          "@type": "HowToSection",
          "name": "Assemble the pie",
        }]
    }
    `)

	r := Recipe{}
	err := r.Read_jsonld(b)
	if err == nil {
		t.Errorf("Expected to Error due to missing 'itemListElement' in HowToSection")
	}
}

// test if HowToStep is missing
func TestBadHowToStep(t *testing.T) {
	b := []byte(`
    {
        "recipeInstructions": 
        [{
          "@type": "HowToStep",
        }]
    }
    `)

	r := Recipe{}
	err := r.Read_jsonld(b)
	if err == nil {
		t.Errorf("Expected to Error due to missing 'Text' in HowToStep")
	}
}

func TestComplexParse(t *testing.T) {
	b := []byte(`
{
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "Article",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#article",
      "isPartOf": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/"
      },
      "author": {
        "name": "Sally",
        "@id": "https://sallysbakingaddiction.com/#/schema/person/3744271f85efc7e271411c4968f1c4ba"
      },
      "headline": "Favorite Apple Pie Recipe",
      "datePublished": "2023-10-05T10:00:31+00:00",
      "dateModified": "2023-06-20T15:14:02+00:00",
      "wordCount": 1810,
      "commentCount": 320,
      "publisher": {
        "@id": "https://sallysbakingaddiction.com/#organization"
      },
      "image": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#primaryimage"
      },
      "thumbnailUrl": "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg",
      "keywords": [
        "all-purpose flour",
        "allspice",
        "apples",
        "cinnamon",
        "eggs",
        "granulated sugar",
        "lemon juice",
        "lemons",
        "nutmeg"
      ],
      "articleSection": [
        "Desserts",
        "Fall",
        "Fruit Pies",
        "Nut Free",
        "Pies, Crisps, &amp; Tarts",
        "Seasonal",
        "Thanksgiving",
        "Videos"
      ],
      "inLanguage": "en-US",
      "potentialAction": [
        {
          "@type": "CommentAction",
          "name": "Comment",
          "target": [
            "https://sallysbakingaddiction.com/apple-pie-recipe/#respond"
          ]
        }
      ]
    },
    {
      "@type": [
        "WebPage",
        "FAQPage"
      ],
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/",
      "url": "https://sallysbakingaddiction.com/apple-pie-recipe/",
      "name": "Favorite Apple Pie Recipe (VIDEO) - Sally&#039;s Baking Addiction",
      "isPartOf": {
        "@id": "https://sallysbakingaddiction.com/#website"
      },
      "primaryImageOfPage": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#primaryimage"
      },
      "image": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#primaryimage"
      },
      "thumbnailUrl": "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg",
      "datePublished": "2023-10-05T10:00:31+00:00",
      "dateModified": "2023-06-20T15:14:02+00:00",
      "description": "With a mountain of gooey cinnamon apples nestled under a flaky pie crust, this is most certainly my favorite apple pie recipe.",
      "breadcrumb": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#breadcrumb"
      },
      "mainEntity": [
        {
          "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501312727"
        },
        {
          "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501419817"
        },
        {
          "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696502837391"
        }
      ],
      "inLanguage": "en-US",
      "potentialAction": [
        {
          "@type": "ReadAction",
          "target": [
            "https://sallysbakingaddiction.com/apple-pie-recipe/"
          ]
        }
      ]
    },
    {
      "@type": "ImageObject",
      "inLanguage": "en-US",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#primaryimage",
      "url": "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg",
      "contentUrl": "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg",
      "width": 1200,
      "height": 1800,
      "caption": "apple pie slice on plate with melty vanilla ice cream scoop on top."
    },
    {
      "@type": "BreadcrumbList",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#breadcrumb",
      "itemListElement": [
        {
          "@type": "ListItem",
          "position": 1,
          "name": "Home",
          "item": "https://sallysbakingaddiction.com/"
        },
        {
          "@type": "ListItem",
          "position": 2,
          "name": "Recipes",
          "item": "https://sallysbakingaddiction.com/recipe-index/"
        },
        {
          "@type": "ListItem",
          "position": 3,
          "name": "Pies, Crisps, &amp; Tarts",
          "item": "https://sallysbakingaddiction.com/category/pies-crisps-tarts"
        },
        {
          "@type": "ListItem",
          "position": 4,
          "name": "Favorite Apple Pie Recipe"
        }
      ]
    },
    {
      "@type": "WebSite",
      "@id": "https://sallysbakingaddiction.com/#website",
      "url": "https://sallysbakingaddiction.com/",
      "name": "Sally&#039;s Baking Addiction",
      "description": "Trusted Recipes from a Self-Taught Baker",
      "publisher": {
        "@id": "https://sallysbakingaddiction.com/#organization"
      },
      "potentialAction": [
        {
          "@type": "SearchAction",
          "target": {
            "@type": "EntryPoint",
            "urlTemplate": "https://sallysbakingaddiction.com/?s={search_term_string}"
          },
          "query-input": "required name=search_term_string"
        }
      ],
      "inLanguage": "en-US"
    },
    {
      "@type": "Organization",
      "@id": "https://sallysbakingaddiction.com/#organization",
      "name": "Sallys Baking Addiction",
      "url": "https://sallysbakingaddiction.com/",
      "logo": {
        "@type": "ImageObject",
        "inLanguage": "en-US",
        "@id": "https://sallysbakingaddiction.com/#/schema/logo/image/",
        "url": "https://sallysbakingaddiction.com/wp-content/uploads/2022/01/android-chrome-512x512-1.png",
        "contentUrl": "https://sallysbakingaddiction.com/wp-content/uploads/2022/01/android-chrome-512x512-1.png",
        "width": 512,
        "height": 512,
        "caption": "Sallys Baking Addiction"
      },
      "image": {
        "@id": "https://sallysbakingaddiction.com/#/schema/logo/image/"
      }
    },
    {
      "@type": "Person",
      "@id": "https://sallysbakingaddiction.com/#/schema/person/3744271f85efc7e271411c4968f1c4ba",
      "name": "Sally",
      "description": "Sally McKenney is a professional food photographer, cookbook author, and baker. Her kitchen-tested recipes and thorough step-by-step tutorials give readers the knowledge and confidence to bake from scratch. Sally has been featured on Good Morning America, HuffPost, Taste of Home, People and more.",
      "sameAs": [
        "https://sallysbakingaddiction.com/about/"
      ]
    },
    {
      "@type": "Question",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501312727",
      "position": 1,
      "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501312727",
      "name": "Should you cook the apples before baking apple pie?",
      "answerCount": 1,
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "You don't <em>have</em> to pre-cook the filling before spooning it into the pie crust, but it's a quick step I recommend. Just 5 minutes on the stove begins the softening process, and also helps the flavors start to mingle. I've never regretted taking this step, and it's certainly catapulted my apple pies from <em>good</em> to <em>great</em>.",
        "inLanguage": "en-US"
      },
      "inLanguage": "en-US"
    },
    {
      "@type": "Question",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501419817",
      "position": 2,
      "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696501419817",
      "name": "How do you make apple pie so the bottom crust isn't soggy?",
      "answerCount": 1,
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "You don't have to pre-bake the bottom pie crust for this pie. There's simply no need to take this extra step because the apple pie bakes for a really long time in the oven. If your pies have soggy crusts, you may not be baking them long enough. See <em>How Do I Know When Apple Pie Is Done</em> above. Additionally, and this is important, I strongly recommend using a <a href=\"https://amzn.to/2LJYp1w\" target=\"_blank\" rel=\"noreferrer noopener sponsored nofollow\">glass pie dish</a>. Glass conducts heat slowly and evenly, and you can literally *see* if the bottom crust is done.",
        "inLanguage": "en-US"
      },
      "inLanguage": "en-US"
    },
    {
      "@type": "Question",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696502837391",
      "position": 3,
      "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#faq-question-1696502837391",
      "name": "What if I don't want to mess with pie crust?",
      "answerCount": 1,
      "acceptedAnswer": {
        "@type": "Answer",
        "text": "Skip traditional pie crust and make my <a href=\"https://sallysbakingaddiction.com/salted-caramel-apple-pie-bars/\">salted caramel apple pie bars</a>, <a href=\"https://sallysbakingaddiction.com/caramel-apple-cheesecake-pie/\">caramel apple cheesecake pie</a>, or my classic <a href=\"https://sallysbakingaddiction.com/apple-crisp/\">apple crisp recipe</a> instead. You might also enjoy my <a href=\"https://sallysbakingaddiction.com/baked-apples/\">baked apples</a>!",
        "inLanguage": "en-US"
      },
      "inLanguage": "en-US"
    },
    {
      "@context": "https://schema.org/",
      "@type": "Recipe",
      "name": "My Best Apple Pie Recipe",
      "description": "With a mountain of gooey, cinnamon-kissed apples nestled under a perfectly buttery and flaky pie crust, this is most certainly my favorite apple pie recipe. To bring out the best apple flavor and texture, pre-cook the filling for only about 5 minutes on the stove. Bake and cool the pie, and then top with vanilla ice cream for the ultimate homestyle dessert.",
      "author": {
        "@type": "Person",
        "name": "Sally",
        "url": "https://sallysbakingaddiction.com/about/"
      },
      "keywords": "apple pie",
      "image": [
        "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2-225x225.jpg",
        "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2-260x195.jpg",
        "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2-320x180.jpg",
        "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg"
      ],
      "url": "https://sallysbakingaddiction.com/apple-pie-recipe/",
      "recipeIngredient": [
        "Homemade Pie Crust or All Butter Pie Crust (both recipes make 2 crusts, 1 for bottom and 1 for top)",
        "10 cups (1250g) 1/4-inch-thick apple slices (about 8 large peeled and cored apples)*",
        "1/2 cup (100g) granulated sugar (or packed brown sugar)",
        "1/4 cup (31g) all-purpose flour (spooned &amp; leveled)",
        "1 Tablespoon (15ml) lemon juice",
        "1 and 1/2 teaspoons ground cinnamon",
        "1/4 teaspoon each: ground allspice &amp; ground nutmeg",
        "egg wash: 1 large egg beaten with 1 Tablespoon (15ml) milk",
        "optional: coarse sugar for sprinkling on crust"
      ],
      "recipeInstructions": [
        {
          "@type": "HowToStep",
          "text": "Prepare either pie crust recipe through step 5.",
          "name": "The crust",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-1"
        },
        {
          "@type": "HowToStep",
          "text": "In a large bowl, stir the apple slices, sugar, flour, lemon juice, cinnamon, allspice, and nutmeg together until thoroughly combined.",
          "name": "Make the filling",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-2"
        },
        {
          "@type": "HowToStep",
          "text": "Pour the apple filling into a very large skillet, or dutch oven, and place over medium-low heat. Stir and cook for 5 minutes until the apples begin to soften. Remove from heat and set aside. This step is optional, but I&#8217;ve found it makes for a juicier, more flavorful filling because it helps begin to soften the apples. If you can, take the few extra minutes to do this, because the flavor is worth it!",
          "name": "Optional pre-cook",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-3"
        },
        {
          "@type": "HowToStep",
          "text": "Preheat oven to 400°F (204°C).",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-4"
        },
        {
          "@type": "HowToStep",
          "text": "On a floured work surface, roll out one of the discs of chilled dough (keep the other one in the refrigerator). Turn the dough about a quarter turn after every few rolls until you have a circle 12 inches in diameter. Carefully place the dough into a 9-inch pie dish that&#8217;s 1.5 to 2 inches deep. Tuck the dough in with your fingers, making sure it is smooth. Spoon the filling into the crust. It&#8217;s ok if it is still warm from the precooking step. It will seem like a lot of apples; that&#8217;s ok. Pile them high, and tightly together.",
          "name": "Roll out the chilled pie dough",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-5"
        },
        {
          "@type": "HowToStep",
          "text": "Remove the other disc of chilled pie dough from the refrigerator. Roll the dough into a circle that is 12 inches diameter. Using a pastry wheel, sharp knife, or pizza cutter, cut strips of dough; in the pictured pie, I cut 12 1-inch-wide strips. Carefully thread the strips over and under one another, pulling back strips as necessary to weave. (Here&#8217;s a lattice pie crust tutorial if you need visuals.) Use a small paring knife or kitchen shears to trim off excess dough. Fold the overhang back towards the center of the pie, and pinch the edges to adhere the top and bottom crusts together. Crimp or flute the pie crust edges to seal.",
          "name": "Finish assembling",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-6"
        },
        {
          "@type": "HowToStep",
          "text": "Lightly brush the top of the pie crust with the egg wash. Sprinkle the top with coarse sugar, if using.",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-7"
        },
        {
          "@type": "HowToStep",
          "text": "Place the pie onto a large baking sheet and bake for 25 minutes. Then, keeping the pie in the oven, reduce the oven temperature down to 375°F (190°C). Place a pie crust shield (see Note for homemade shield) on the edges to prevent them from over-browning. Continue baking the pie until the filling is bubbling around the edges, 35–40 more minutes. This sounds like a long time, but under-baking the pie means an unfinished filling with firm apples with paste-like flour. If you want to be precise, the internal temperature of the filling taken with an instant read thermometer should be around 200°F (93°C) when done. Tip: If needed towards the end of bake time, remove the pie crust shield and tent an entire piece of foil on top of the pie if the top looks like it&#8217;s getting too brown.",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-8"
        },
        {
          "@type": "HowToStep",
          "text": "Remove pie from the oven, place on a cooling rack, and cool for at least 3 hours before slicing and serving. Filling will be too juicy if the pie is warm when you slice it.",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-9"
        },
        {
          "@type": "HowToStep",
          "text": "​​Cover and store leftover pie at room temperature for up to 1 day or in the refrigerator for up to 5 days.",
          "url": "https://sallysbakingaddiction.com/apple-pie-recipe/#instruction-step-10"
        }
      ],
      "prepTime": "PT3H",
      "cookTime": "PT1H5M",
      "totalTime": "PT7H",
      "recipeYield": [
        "8",
        "8-10 servings"
      ],
      "recipeCategory": "Pie",
      "cookingMethod": "Baking",
      "recipeCuisine": "American",
      "aggregateRating": {
        "@type": "AggregateRating",
        "reviewCount": "76",
        "ratingValue": "4.9"
      },
      "video": {
        "@context": "http://schema.org",
        "@type": "VideoObject",
        "name": "Apple Pie | Sally's Baking Recipes",
        "description": "This deep-dish apple pie recipe features layers upon layers of sweet spiced apples nestled in a buttery flaky pie crust. \n\nGet the full recipe: https://sallysbakingaddiction.com/deep-dish-apple-pie/\n\n• Ask your recipe question or leave a review over on the recipe page.\n\n#baking #recipes #applepie \n\n• More of Sally's baking recipes: https://sallysbakingaddiction.com/",
        "duration": "PT3M",
        "embedUrl": "https://www.youtube.com/embed/keyJHNHimOM?feature=oembed",
        "contentUrl": "https://www.youtube.com/watch?v=keyJHNHimOM",
        "thumbnailUrl": [
          "https://i.ytimg.com/vi/keyJHNHimOM/hqdefault.jpg"
        ],
        "uploadDate": "2023-06-15T13:26:45+00:00"
      },
      "review": [
        {
          "@type": "Review",
          "reviewRating": {
            "@type": "Rating",
            "ratingValue": "5"
          },
          "author": {
            "@type": "Person",
            "name": "Emily"
          },
          "datePublished": "2024-03-19",
          "reviewBody": "I've been making this pie every weekend because it's delicious and I want to perfect it! Really appreciate the detailed instructions on the pie and your shortening crust. I've been pre-cooking the filling as instructed, but I end up with clumps of cooked filling (flour+spices+sugar) that remain in clumps in the finished pie. Any ideas for how to prevent this?"
        },
        {
          "@type": "Review",
          "reviewRating": {
            "@type": "Rating",
            "ratingValue": "5"
          },
          "author": {
            "@type": "Person",
            "name": "XY"
          },
          "datePublished": "2024-05-13",
          "reviewBody": "I halved the recipe and it smelt so good I had it at 10pm, and I don’t regret it. The all butter crust is difficult to handle but totally worth it. And I would follow her recommendation on pre cooking the apples to give them a slight head start."
        }
      ],
      "datePublished": "2023-10-05",
      "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#recipe",
      "isPartOf": {
        "@id": "https://sallysbakingaddiction.com/apple-pie-recipe/#article"
      },
      "mainEntityOfPage": "https://sallysbakingaddiction.com/apple-pie-recipe/"
    }
  ]
}
`)

	r := Recipe{}
	err := r.Read_jsonld(b)
	if err != nil {
		t.Errorf("%s", err)
	}

	if r.Name != "My Best Apple Pie Recipe" {
		t.Errorf("Recipe Name incorrect. Expected: My Best Apple Pie Recipe, Actual: %s", r.Name)
	}

	if r.ImageUrl != "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg" {
		t.Errorf("Recipe Image Url incorrect. Expected: https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2.jpg, Actual: %s", r.ImageUrl)
	}
     
	if r.ThumbnailUrl != "https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2-260x195.jpg" {
		t.Errorf("Recipe Thumbnail Url incorrect. Expected: https://sallysbakingaddiction.com/wp-content/uploads/2017/07/slice-of-apple-pie-2-260x195.jpg, Actual: %s", r.ThumbnailUrl)
	}

	ingredients := []string{
		"Homemade Pie Crust or All Butter Pie Crust (both recipes make 2 crusts, 1 for bottom and 1 for top)",
		"10 cups (1250g) 1/4-inch-thick apple slices (about 8 large peeled and cored apples)*",
		"1/2 cup (100g) granulated sugar (or packed brown sugar)",
		"1/4 cup (31g) all-purpose flour (spooned &amp; leveled)",
		"1 Tablespoon (15ml) lemon juice",
		"1 and 1/2 teaspoons ground cinnamon",
		"1/4 teaspoon each: ground allspice &amp; ground nutmeg",
		"egg wash: 1 large egg beaten with 1 Tablespoon (15ml) milk",
		"optional: coarse sugar for sprinkling on crust",
	}
	ingredientsMatch(t, ingredients, r)

	instructions := []string{
		"Prepare either pie crust recipe through step 5.",
		"In a large bowl, stir the apple slices, sugar, flour, lemon juice, cinnamon, allspice, and nutmeg together until thoroughly combined.",
		"Pour the apple filling into a very large skillet, or dutch oven, and place over medium-low heat. Stir and cook for 5 minutes until the apples begin to soften. Remove from heat and set aside. This step is optional, but I&#8217;ve found it makes for a juicier, more flavorful filling because it helps begin to soften the apples. If you can, take the few extra minutes to do this, because the flavor is worth it!",
		"Preheat oven to 400°F (204°C).",
		"On a floured work surface, roll out one of the discs of chilled dough (keep the other one in the refrigerator). Turn the dough about a quarter turn after every few rolls until you have a circle 12 inches in diameter. Carefully place the dough into a 9-inch pie dish that&#8217;s 1.5 to 2 inches deep. Tuck the dough in with your fingers, making sure it is smooth. Spoon the filling into the crust. It&#8217;s ok if it is still warm from the precooking step. It will seem like a lot of apples; that&#8217;s ok. Pile them high, and tightly together.",
		"Remove the other disc of chilled pie dough from the refrigerator. Roll the dough into a circle that is 12 inches diameter. Using a pastry wheel, sharp knife, or pizza cutter, cut strips of dough; in the pictured pie, I cut 12 1-inch-wide strips. Carefully thread the strips over and under one another, pulling back strips as necessary to weave. (Here&#8217;s a lattice pie crust tutorial if you need visuals.) Use a small paring knife or kitchen shears to trim off excess dough. Fold the overhang back towards the center of the pie, and pinch the edges to adhere the top and bottom crusts together. Crimp or flute the pie crust edges to seal.",
		"Lightly brush the top of the pie crust with the egg wash. Sprinkle the top with coarse sugar, if using.",
		"Place the pie onto a large baking sheet and bake for 25 minutes. Then, keeping the pie in the oven, reduce the oven temperature down to 375°F (190°C). Place a pie crust shield (see Note for homemade shield) on the edges to prevent them from over-browning. Continue baking the pie until the filling is bubbling around the edges, 35–40 more minutes. This sounds like a long time, but under-baking the pie means an unfinished filling with firm apples with paste-like flour. If you want to be precise, the internal temperature of the filling taken with an instant read thermometer should be around 200°F (93°C) when done. Tip: If needed towards the end of bake time, remove the pie crust shield and tent an entire piece of foil on top of the pie if the top looks like it&#8217;s getting too brown.",
		"Remove pie from the oven, place on a cooling rack, and cool for at least 3 hours before slicing and serving. Filling will be too juicy if the pie is warm when you slice it.",
		"​​Cover and store leftover pie at room temperature for up to 1 day or in the refrigerator for up to 5 days.",
	}

	instructionsMatch(t, instructions, r)
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
