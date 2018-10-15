package main

import (
	"testing"
)

func TestEqual(t *testing.T) {
	drink1 := Drink{
		Title:    "Bellini",
		Category: "Sparkling Cocktail",
		Ingredients: []string{
			"10 cl Prosecco",
			"5 cl Fresh peach puree",
		},
		Instructions: "Pour peach puree into chilled glass and add sparkling wine. Stir gently.",
		Link:         "http://iba-world.com/cocktails/bellini/",
	}

	drink2 := Drink{
		Title:    "Bellini",
		Category: "Sparkling Cocktail",
		Ingredients: []string{
			"10 cl Prosecco",
			"5 cl Fresh peach puree",
		},
		Instructions: "Pour peach puree into chilled glass and add sparkling wine. Stir gently.",
		Link:         "http://iba-world.com/cocktails/bellini/",
	}

	if !drink1.Equal(drink2) {
		t.Error()
	}

}

func TestGetDrinkByName(t *testing.T) {
	want := Drink{
		Title:    "Bellini",
		Category: "Sparkling Cocktail",
		Ingredients: []string{
			"10 cl Prosecco",
			"5 cl Fresh peach puree",
		},
		Instructions: "Pour peach puree into chilled glass and add sparkling wine. Stir gently.",
		Link:         "http://iba-world.com/cocktails/bellini/",
	}
	got, err := getDrinkByName("bellini")
	if err != nil {
		t.Error("Got err from getDrinkByName")
	}
	if !want.Equal(got) {
		t.Error("Drinks deemed not equal")
	}
}
